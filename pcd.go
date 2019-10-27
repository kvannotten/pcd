// Copyright © 2018 Kristof Vannotten <kristof@vannotten.be>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package pcd

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	urlpath "path"
	"path/filepath"
	"strings"

	"github.com/kvannotten/pcd/rss"
	"github.com/pkg/errors"
)

type Podcast struct {
	ID   int
	Name string
	Feed string
	Path string

	// Login data if there's authentication involved
	Username string
	Password string

	// List of episodes
	Episodes []Episode
}

type Episode struct {
	Title string
	Date  string
	URL   string
}

var (
	ErrCouldNotSync          = errors.New("Could not sync podcast")
	ErrRequestFailed         = errors.New("Could not perform request")
	ErrAccessDenied          = errors.New("Access denied to feed")
	ErrFilesystemError       = errors.New("Could not do filesystem request")
	ErrParserIssue           = errors.New("Could not parse feed")
	ErrEncodeError           = errors.New("Could not encode feed")
	ErrFeedNotFound          = errors.New("Could not find feed (404)")
	ErrCouldNotDownload      = errors.New("Could not download episode")
	ErrCouldNotReadFromCache = errors.New("Could not read episodes from cache. Perform a sync and try again.")
	ErrCouldNotParseContent  = errors.New("Could not parse the content from the feed")
)

func (p *Podcast) Sync() error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", p.Feed, nil)
	if err != nil {
		log.Print(err)
		return ErrCouldNotSync
	}

	if p.Username != "" {
		req.SetBasicAuth(p.Username, p.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return ErrRequestFailed
	}
	switch resp.StatusCode {
	case http.StatusOK: // NOOP
	case http.StatusForbidden, http.StatusUnauthorized:
		return ErrAccessDenied
	case http.StatusNotFound:
		return ErrFeedNotFound
	case http.StatusInternalServerError:
		return ErrRequestFailed
	default:
		return ErrRequestFailed
	}
	defer resp.Body.Close()

	p.Episodes, err = parseEpisodes(resp.Body)
	if err != nil {
		log.Print(err)
		return ErrParserIssue
	}

	if err := os.MkdirAll(p.Path, os.ModePerm); err != nil {
		log.Print(err)
		return ErrFilesystemError
	}

	path := filepath.Join(p.Path, ".feed")
	f, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return ErrFilesystemError
	}
	defer f.Close()

	blob, err := toGOB64(p.Episodes)
	if err != nil {
		log.Print(err)
		return ErrEncodeError
	}
	if _, err := io.Copy(f, blob); err != nil {
		log.Print(err)
		return ErrFilesystemError
	}

	return nil
}

func (p *Podcast) Load() error {
	path := filepath.Join(p.Path, ".feed")
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Could not open feed file: %#v", err)
		return ErrCouldNotReadFromCache
	}
	defer f.Close()

	p.Episodes, err = fromGOB64(f)
	if err != nil {
		log.Printf("Could not decode episodes: %#v", err)
		return ErrCouldNotReadFromCache
	}

	return nil
}

const (
	titleLength = 60
)

func (p *Podcast) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("All episodes of %s (id: %d)\n", p.Name, p.ID))

	// find longest episode title to see if title length is smaller than titleLength
	tl := 0
	for _, episode := range p.Episodes {
		if len(episode.Title) > tl {
			tl = len(episode.Title)
		}
	}
	if tl > titleLength {
		tl = titleLength
	}

	for index, episode := range p.Episodes {
		title := episode.Title
		if len(episode.Title) > titleLength {
			title = fmt.Sprintf("%s...", episode.Title[0:(titleLength-4)])
		}
		formatStr := fmt.Sprintf("%%-4d %%-%ds %%20s\n", tl)
		sb.WriteString(fmt.Sprintf(formatStr, index+1, title, episode.Date))
	}

	return sb.String()
}

// Download downloads an episode in 'path'. The writer argument is optional
// and will just mirror everything written into it (useful for tracking the speed)
func (e *Episode) Download(path string, writer io.Writer) error {
	u, err := url.Parse(e.URL)
	if err != nil {
		log.Printf("Parse episode url failed: %#v", err)
		return ErrCouldNotDownload
	}

	if u.Path == "" {
		return ErrFilesystemError
	}

	// remove the query string from filename
	q := u.Query()
	for k := range q {
		q.Del(k)
	}

	filename := urlpath.Base(u.Path)
	fpath := filepath.Join(path, filename)

	if _, err := os.Stat(fpath); !os.IsNotExist(err) {
		return ErrFilesystemError
	}

	res, err := http.Get(e.URL)
	if err != nil {
		log.Printf("Could not download episode: %#v", err)
		return ErrCouldNotDownload
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("Could not download episode: %#v", err)
		return ErrCouldNotDownload
	}

	f, err := os.Create(fpath)
	if err != nil {
		log.Printf("Could not create file: %#v", err)
		return ErrCouldNotDownload
	}
	defer f.Close()

	var mw io.Writer

	if writer != nil {
		mw = io.MultiWriter(f, writer)
	} else {
		mw = f
	}
	if _, err := io.Copy(mw, res.Body); err != nil {
		log.Printf("Could not write to file: %#v", err)
		return ErrCouldNotDownload
	}

	return nil
}

func parseEpisodes(content io.Reader) ([]Episode, error) {
	feed, err := rss.Parse(content)
	if err != nil {
		return nil, ErrCouldNotParseContent
	}

	var episodes []Episode

	for _, item := range feed.Channel.Items {

		episode := Episode{
			Title: item.Title.Title,
			Date:  item.Date.Date,
			URL:   item.Enclosure.URL,
		}

		episodes = append(episodes, episode)
	}

	return episodes, nil
}

func toGOB64(episodes []Episode) (io.Reader, error) {
	b := bytes.Buffer{}

	e := gob.NewEncoder(&b)
	if err := e.Encode(episodes); err != nil {
		return nil, err
	}

	dst := bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, &dst)
	encoder.Write(b.Bytes())

	defer encoder.Close()

	return &dst, nil
}

func fromGOB64(content io.Reader) ([]Episode, error) {
	var episodes []Episode

	decoder := base64.NewDecoder(base64.StdEncoding, content)
	d := gob.NewDecoder(decoder)

	if err := d.Decode(&episodes); err != nil {
		return nil, err
	}

	return episodes, nil
}
