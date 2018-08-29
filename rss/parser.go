package rss

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
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

const Layout string = "Mon, _2 Jan 2006 15:04:05 -0700"

func sortFeedByDate(feed *PodcastFeed) {
	if len(feed.Channel.Items) < 1 {
		return
	}

	firstDate, _ := time.Parse(Layout, feed.Channel.Items[0].Date.Date)
	lastDate, _ := time.Parse(Layout, feed.Channel.Items[len(feed.Channel.Items)-1].Date.Date)

	if firstDate.After(lastDate) {
		// reverse the feed
		for i, j := 0, len(feed.Channel.Items)-1; i < j; i, j = i+1, j-1 {
			feed.Channel.Items[i], feed.Channel.Items[j] = feed.Channel.Items[j], feed.Channel.Items[i]
		}
	}
}
