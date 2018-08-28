package pcd

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	Title  string
	Date   time.Time
	URL    string
	Length int
}

var (
	ErrCouldNotSync    = errors.New("Could not sync podcast")
	ErrRequestFailed   = errors.New("Could not perform request")
	ErrAccessDenied    = errors.New("Access denied to feed")
	ErrFilesystemError = errors.New("Could not do filesystem request")
	ErrParserIssue     = errors.New("Could not parse feed")
	ErrEncodeError     = errors.New("Could not encode feed")
	ErrFeedNotFound    = errors.New("Could not find feed (404)")
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

	episodes, err := parseEpisodes(resp.Body)
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

	blob, err := toGOB64(episodes)
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

func parseEpisodes(content io.Reader) ([]Episode, error) {
	feed, err := rss.Parse(content)
	if err != nil {
		return nil, err
	}

	var episodes []Episode

	for _, item := range feed.Channel.Items {
		t, err := time.Parse(rss.Layout, item.Date.Date)
		if err != nil {
			log.Printf("Could not parse episode: %#v", err)
			continue
		}
		episode := Episode{
			Title:  item.Title.Title,
			Date:   t,
			URL:    item.Enclosure.URL,
			Length: item.Enclosure.Length,
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
