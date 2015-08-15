package feedparser

import (
	"bufio"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"github.com/kvannotten/pcd/configuration"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type PodcastFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       ChannelTitle
	Description ChannelDescription
	Items       []Item `xml:"item"`
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
	Title string
	//Link      string   `xml:",chardata"`
	//Guid      string   `xml:",chardata"`
	Enclosure Enclosure
}

type ItemTitle struct {
	XMLName xml.Name `xml:"title"`
	Title   string   `xml:",chardata"`
}

type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
	Length  int64    `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

func Parse(podcast configuration.Podcast) {
	resp, err := http.Get(podcast.Feed)
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

func Download(podcast configuration.Podcast) {
	feed := readCachedFeed(podcast)
	fmt.Println(feed)
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
