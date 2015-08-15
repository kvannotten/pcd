package feedparser

import (
	"bufio"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvannotten/pcd/configuration"
)

type PodcastFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Items       []Item   `xml:"item"`
	Title       ChannelTitle
	Description ChannelDescription
}

type ChannelTitle struct {
	XMLName xml.Name `xml:"title"`
	Title   string   `xml:",chardata"`
}

type ChannelDescription struct {
	XMLName     xml.Name `xml:"description"`
	Description string   `xml:",chardata"`
}

type Item struct {
	Title      ItemTitle
	Enclosure  Enclosure
	Downloaded bool
}

type ItemTitle struct {
	XMLName xml.Name `xml:"title"`
	Title   string   `xml:",chardata"`
}

type ItemLink struct {
	XMLName xml.Name `xml:"link"`
	Link    string   `xml:",chardata"`
}

type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  uint64   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

func Parse(podcast configuration.Podcast) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", podcast.Feed, nil)
	if err != nil {
		fmt.Printf("Could not build request: %s\n", err)
		return
	}

	if len(podcast.Username) > 0 && len(podcast.Password) > 0 {
		req.SetBasicAuth(podcast.Username, podcast.Password)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	var feed PodcastFeed
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err := xml.Unmarshal(body, &feed); err != nil {
		return
	}
	writeFeed(podcast, feed)
}

func Download(podcast configuration.Podcast, number int) {
	feed := readCachedFeed(podcast)
	url := feed.Channel.Items[number-1].Enclosure.URL

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if len(podcast.Username) > 0 && len(podcast.Password) > 0 {
		req.SetBasicAuth(podcast.Username, podcast.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Could not download podcast: %s\n", err)
		return
	}
	defer resp.Body.Close()

	tokens := strings.Split(url, "/")
	filename := tokens[len(tokens)-1]

	writePodcast(podcast, resp.Body, filename)
}

func ListEpisodes(podcast configuration.Podcast) []Item {
	items := make([]Item, 0)
	feed := readCachedFeed(podcast)

	for i := 0; i < len(feed.Channel.Items); i++ {
		item := feed.Channel.Items[i]

		tokens := strings.Split(item.Enclosure.URL, "/")
		filename := tokens[len(tokens)-1]
		path := filepath.Join(podcast.Path, filename)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			item.Downloaded = false
		} else {
			item.Downloaded = true
		}
		items = append(items, item)
	}

	return items
}

func writePodcast(podcast configuration.Podcast, reader io.Reader, filename string) {
	path := filepath.Join(podcast.Path, filename)
	fmt.Printf("Downloading podcast to %s\n", path)

	f, err := os.Create(path)
	if err != nil {
		panic("Could not create file")
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	if err != nil {
		panic("Could not download file")
	}
}

func writeFeed(podcast configuration.Podcast, feed PodcastFeed) {
	err := os.MkdirAll(podcast.Path, 0700)
	path := filepath.Join(podcast.Path, ".cache")
	f, err := os.Create(path)
	if err != nil {

	}
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	enc.Encode(feed)

	w.Flush()
}

func readCachedFeed(podcast configuration.Podcast) PodcastFeed {
	path := filepath.Join(podcast.Path, ".cache")
	f, err := os.Open(path)
	if err != nil {

	}
	defer f.Close()
	r := bufio.NewReader(f)
	dec := gob.NewDecoder(r)
	var feed PodcastFeed
	dec.Decode(&feed)

	return feed
}
