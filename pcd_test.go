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
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSync(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	podcast := &Podcast{
		ID:   1,
		Name: "test",
		Feed: ts.URL,
		Path: randomPath(t),
	}

	if err := podcast.Sync(); err != nil {
		t.Errorf("Expected to be able to sync, but could not sync: %#v", err)
	}
}

func TestSyncBadRequest(t *testing.T) {
	podcast := &Podcast{
		ID:   1,
		Name: "test",
		Feed: "foo",
	}

	if err := podcast.Sync(); err != ErrRequestFailed {
		t.Errorf("Expected %#v, but got: %#v", ErrRequestFailed, err)
	}
}

func TestSyncPathIssue(t *testing.T) {
	// test is running as root and will fail
	// so just skip it
	if os.Geteuid() == 0 {
		t.Skip()
	}

	ts := testServer()
	defer ts.Close()

	podcast := &Podcast{
		ID:   1,
		Name: "test",
		Feed: ts.URL,
		Path: "/root/access/required",
	}

	if err := podcast.Sync(); err != ErrFilesystemError {
		t.Errorf("Expected %#v, but got: %#v", ErrFilesystemError, err)
	}
}

func TestCredentials(t *testing.T) {
	ts := testServerWithBasicAuth("test", "foo")

	podcast := &Podcast{
		ID:       1,
		Name:     "test",
		Feed:     ts.URL,
		Username: "incorrect",
		Password: "incorrect",
	}

	if err := podcast.Sync(); err != ErrAccessDenied {
		t.Errorf("Expected %#v, but got: %#v", ErrAccessDenied, err)
	}
}

func TestFeedRequestIssues(t *testing.T) {
	table := []struct {
		name string
		got  int
		want error
	}{
		{"feed not found", http.StatusNotFound, ErrFeedNotFound},
		{"issue with server", http.StatusInternalServerError, ErrRequestFailed},
		{"random non-200", http.StatusBadRequest, ErrRequestFailed},
	}
	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			ts := testServerWithStatusCode(e.got)

			podcast := &Podcast{
				Feed: ts.URL,
			}

			if err := podcast.Sync(); err != e.want {
				t.Errorf("Expected %#v, but got: %#v", e.want, err)
			}

			ts.Close()
		})
	}
}

func TestEpisodes(t *testing.T) {
	ts := testServer()

	podcast := &Podcast{
		Name: "test",
		Feed: ts.URL,
		Path: randomPath(t),
	}

	if err := podcast.Sync(); err != nil {
		t.Errorf("Expected no error, but got: %#v", err)
	}

	if len(podcast.Episodes) != 1 {
		t.Errorf("Expected to have 1 episode, but got: %d", len(podcast.Episodes))
	}

	if podcast.Episodes[0].Title != "Title of Podcast Episode" {
		t.Errorf("Expected episode to have title 'Title of Podcast Episode', but got: %s", podcast.Episodes[0].Title)
	}
}

func TestPodcastString(t *testing.T) {
	now := "Thu, 10 Nov 2016 19:41:48 -0700"

	podcast := &Podcast{
		Name:     "test",
		Episodes: []Episode{{Title: "foo", Date: now}},
	}

	if !strings.Contains(podcast.String(), podcast.Name) {
		t.Error("Expected podcast name to be in the output")
	}
	if !strings.Contains(podcast.String(), podcast.Episodes[0].Title) {
		t.Error("Expected episode title to be in the output")
	}
	if !strings.Contains(podcast.String(), podcast.Episodes[0].Date) {
		t.Error("Expected date to be in the output")
	}

	t.Run("extra long name", func(t *testing.T) {
		podcast.Episodes[0].Title = "this is a really long name that should be truncated to something shorter to fit the screen of the user"
		if !strings.Contains(podcast.String(), podcast.Episodes[0].Title[0:titleLength-4]) {
			t.Error("Expected truncated episode title to be in the output")
		}
	})
}

func TestLoad(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	podcast := &Podcast{
		ID:   1,
		Name: "test",
		Feed: ts.URL,
		Path: randomPath(t),
	}

	if err := podcast.Sync(); err != nil {
		t.Errorf("Could not sync podcast: %#v", err)
	}
	podcast.Episodes = nil
	if len(podcast.Episodes) != 0 {
		t.Errorf("Episodes should be empty")
	}

	faultyFeedPath := randomPath(t)
	f, err := os.Create(filepath.Join(faultyFeedPath, ".feed"))
	if err != nil {
		t.Error("Error while creating temporary file...")
	}
	f.WriteString("invalid data")
	f.Close()

	table := []struct {
		name          string
		path          string
		err           error
		checkEpisodes bool
	}{
		{"valid load", podcast.Path, nil, true},
		{"valid path but no .feed file", randomPath(t), ErrCouldNotReadFromCache, false},
		{"invalid path", "/root/access", ErrCouldNotReadFromCache, false},
		{"valid path but faulty .feed file", faultyFeedPath, ErrCouldNotReadFromCache, false},
	}

	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			podcast.Path = e.path

			if err := podcast.Load(); err != e.err {
				t.Errorf("Expected %#v, but got: %#v", e.err, err)
			}

			if e.checkEpisodes && len(podcast.Episodes) != 1 {
				t.Errorf("Expected 1 podcast episode to be present, but got: %d", len(podcast.Episodes))
			}
		})
	}

}

func TestDownload(t *testing.T) {
	r := strings.NewReader(Podcastfeed)
	episodes, err := parseEpisodes(r)
	if err != nil {
		t.Errorf("Expected no error, but got: %#v", err)
	}

	episode := episodes[0]
	ts := testServer()
	episode.URL = ts.URL + "/sample.mp3"
	episode.Title = "some_title"
	episode.Guid = "some_guid"

	if err := episode.Download(randomPath(t), nil); err != nil {
		t.Errorf("Expected to be able to download episode, but got: %#v", err)
	}
}

func TestInvalidDownload(t *testing.T) {
	table := []struct {
		name    string
		episode *Episode
		path    string
		writer  io.Writer
		err     error
	}{
		{"invalid url", &Episode{URL: "invalid"}, randomPath(t), nil, ErrCouldNotDownload},
		{"invalid status", &Episode{URL: testServerWithStatusCode(404).URL + "/sample.mp3"}, randomPath(t), nil, ErrCouldNotDownload},
		{"invalid path", &Episode{URL: testServer().URL}, "/root/access", nil, ErrFilesystemError},
	}

	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			if err := e.episode.Download(e.path, nil); err != e.err {
				t.Errorf("Expected %#v, but got %#v", e.err, err)
			}
		})
	}
}

func TestParseEpisodes(t *testing.T) {
	table := []struct {
		name        string
		feed        io.Reader
		err         error
		hasEpisodes bool
		title       string
	}{
		{"valid feed", strings.NewReader(Podcastfeed), nil, true, "Title of Podcast Episode"},
		{"invalid feed should return error", strings.NewReader("some invalid text"), ErrCouldNotParseContent, false, ""},
		{"invalid episode should just continue", strings.NewReader(invalidEpisodesFeed), nil, false, ""},
	}
	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			episodes, err := parseEpisodes(e.feed)
			if err != e.err {
				t.Errorf("Expected %#v, but got: %#v", e.err, err)
			}
			if e.hasEpisodes {
				if episodes == nil {
					t.Errorf("Expected episodes to not be nil")
				}
				if len(episodes) != 1 {
					t.Errorf("Expected 1 episode, but got: %#v", len(episodes))
				}
				if episodes[0].Title != e.title {
					t.Errorf("Expected title to be '%s', but got %s", e.title, episodes[0].Title)
				}

			}

		})
	}

}

func TestGobEncodeAndDecode(t *testing.T) {
	episode := Episode{
		Title: "test",
		Guid:  "guid",
	}

	content, err := toGOB64([]Episode{episode})
	if err != nil {
		t.Errorf("Didn't expect an error, but got: %#v", err)
	}

	data, err := ioutil.ReadAll(content)
	if err != nil {
		t.Errorf("Didn't expect an error, but got: %#v", err)
	}

	episodes, err := fromGOB64(bytes.NewBuffer(data))
	if err != nil {
		t.Errorf("Didn't expect an error, but got: %#v", err)
	}
	if episodes == nil {
		t.Errorf("Expected episodes to not be nil")
	}
	if len(episodes) != 1 {
		t.Errorf("Expected 1 episode, but got: %#v", len(episodes))
	}
	if episodes[0].Title != episode.Title {
		t.Errorf("Expected title to be %s, but got %s", episode.Title, episodes[0].Title)
	}
}

func TestRepeatableFileName(t *testing.T) {
	episode := Episode{
		Title: "This is a cool episode",
		URL:   "https://some_fake_website.com/episodes/cool.mp3",
		Guid:  "https://some_fake_website.com/blog/cool",
	}

	urlData, err := url.Parse(episode.URL)

	if err != nil {
		t.Errorf("Didn't expect an error, but got: %#v", err)
	}

	if episode.FileName(urlData) != episode.FileName(urlData) {
		t.Errorf("file names should be repeatable")
	}
}
