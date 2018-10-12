package rss

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"time"
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
	Date       PodcastDate
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
	Length  int      `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

type PodcastDate struct {
	XMLName xml.Name `xml:"pubDate"`
	Date    string   `xml:",chardata"`
}

var (
	ErrCouldNotGetContent   = errors.New("Could not get content")
	ErrCouldNotParseContent = errors.New("Could not parse content")
)

func Parse(content io.Reader) (*PodcastFeed, error) {
	if content == nil {
		return nil, ErrCouldNotGetContent
	}

	var feed PodcastFeed
	body, err := ioutil.ReadAll(content)
	if err != nil {
		log.Print(err)
		return nil, ErrCouldNotGetContent
	}
	if err := xml.Unmarshal(body, &feed); err != nil {
		log.Print(err)
		return nil, ErrCouldNotParseContent
	}
	sortFeedByDate(&feed)
	return &feed, nil
}

func stringToDate(d string) time.Time {
	var t time.Time
	var err error

	t, err = time.Parse(time.RFC1123, d)
	if err != nil {
		t, _ = time.Parse(time.RFC1123Z, d)
	}
	return t
}

func sortFeedByDate(feed *PodcastFeed) {
	sort.Slice(feed.Channel.Items, func(i, j int) bool {
		d1 := stringToDate(feed.Channel.Items[i].Date.Date)
		d2 := stringToDate(feed.Channel.Items[j].Date.Date)

		return d2.After(d1)
	})
}
