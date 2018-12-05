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
