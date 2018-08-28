package pcd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
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
		Path: os.TempDir(),
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

func TestParseEpisodes(t *testing.T) {
	r := strings.NewReader(Podcastfeed)

	episodes, err := parseEpisodes(r)
	if err != nil {
		t.Errorf("Didn't expect an error, but got: %#v", err)
	}
	if episodes == nil {
		t.Errorf("Expected episodes to not be nil")
	}
	if len(episodes) != 1 {
		t.Errorf("Expected 1 episode, but got: %#v", len(episodes))
	}
	if episodes[0].Title != "Title of Podcast Episode" {
		t.Errorf("Expected title to be 'Title of Podcast Episode', but got %s", episodes[0].Title)
	}
}

func TestGobEncodeAndDecode(t *testing.T) {
	episode := Episode{
		Title: "test",
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
