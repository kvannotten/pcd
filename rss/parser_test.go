package rss

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

var podcastfeed = `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" version="2.0">
<channel>
<title>Title of Podcast</title>
<link>http://www.example.com/</link>
<language>en-us</language>
<itunes:subtitle>Subtitle of podcast</itunes:subtitle>
<itunes:author>Author Name</itunes:author>
<itunes:summary>Description of podcast.</itunes:summary>
<description>Description of podcast.</description>
<itunes:owner>
    <itunes:name>Owner Name</itunes:name>
    <itunes:email>me@example.com</itunes:email>
</itunes:owner>
<itunes:explicit>no</itunes:explicit>
<itunes:image href="http://www.example.com/podcast-icon.jpg" />
<itunes:category text="Category Name"></itunes:category>

<!--REPEAT THIS BLOCK FOR EACH EPISODE-->
<item>
    <title>Title of Podcast Episode</title>
    <itunes:summary>Description of podcast episode content</itunes:summary>
    <description>Description of podcast episode content</description>
    <link>http://example.com/podcast-1</link>
    <enclosure url="http://example.com/podcast-1/podcast.mp3" type="audio/mpeg" length="1024"></enclosure>
    <pubDate>Thu, 21 Dec 2016 16:01:07 +0000</pubDate>
    <itunes:author>Author Name</itunes:author>
    <itunes:duration>00:32:16</itunes:duration>
    <itunes:explicit>no</itunes:explicit>
    <guid>http://example.com/podcast-1</guid>
</item> 
<item>
    <title>Title of Podcast Episode 2</title>
    <itunes:summary>Description of podcast episode content</itunes:summary>
    <description>Description of podcast episode content</description>
    <link>http://example.com/podcast-1</link>
    <enclosure url="http://example.com/podcast-1/podcast.mp3" type="audio/mpeg" length="1024"></enclosure>
    <pubDate>Thu, 29 Dec 2016 16:01:07 +0000</pubDate>
    <itunes:author>Author Name</itunes:author>
    <itunes:duration>00:32:16</itunes:duration>
    <itunes:explicit>no</itunes:explicit>
    <guid>http://example.com/podcast-1</guid>
</item> 
<!--END REPEAT--> 
   
</channel>
</rss>`

var podcastfeedReversed = `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" version="2.0">
<channel>
<title>Title of Podcast</title>
<link>http://www.example.com/</link>
<language>en-us</language>
<itunes:subtitle>Subtitle of podcast</itunes:subtitle>
<itunes:author>Author Name</itunes:author>
<itunes:summary>Description of podcast.</itunes:summary>
<description>Description of podcast.</description>
<itunes:owner>
    <itunes:name>Owner Name</itunes:name>
    <itunes:email>me@example.com</itunes:email>
</itunes:owner>
<itunes:explicit>no</itunes:explicit>
<itunes:image href="http://www.example.com/podcast-icon.jpg" />
<itunes:category text="Category Name"></itunes:category>

<!--REPEAT THIS BLOCK FOR EACH EPISODE-->
<item>
    <title>Title of Podcast Episode 2</title>
    <itunes:summary>Description of podcast episode content</itunes:summary>
    <description>Description of podcast episode content</description>
    <link>http://example.com/podcast-1</link>
    <enclosure url="http://example.com/podcast-1/podcast.mp3" type="audio/mpeg" length="1024"></enclosure>
    <pubDate>Thu, 29 Dec 2016 16:01:07 +0000</pubDate>
    <itunes:author>Author Name</itunes:author>
    <itunes:duration>00:32:16</itunes:duration>
    <itunes:explicit>no</itunes:explicit>
    <guid>http://example.com/podcast-1</guid>
</item> 
<item>
    <title>Title of Podcast Episode</title>
    <itunes:summary>Description of podcast episode content</itunes:summary>
    <description>Description of podcast episode content</description>
    <link>http://example.com/podcast-1</link>
    <enclosure url="http://example.com/podcast-1/podcast.mp3" type="audio/mpeg" length="1024"></enclosure>
    <pubDate>Thu, 21 Dec 2016 16:01:07 +0000</pubDate>
    <itunes:author>Author Name</itunes:author>
    <itunes:duration>00:32:16</itunes:duration>
    <itunes:explicit>no</itunes:explicit>
    <guid>http://example.com/podcast-1</guid>
</item> 
<!--END REPEAT--> 
   
</channel>
</rss>`

func TestParse(t *testing.T) {
	feed, err := Parse(strings.NewReader(podcastfeed))
	if err != nil {
		t.Errorf("Did not expect error but got: %#v", err)
	}

	table := []struct {
		name string
		got  string
		want string
	}{
		{"podcast title", feed.Channel.Title.Title, "Title of Podcast"},
		{"podcast description", feed.Channel.Description.Description, "Description of podcast."},
		{"title of item", feed.Channel.Items[0].Title.Title, "Title of Podcast Episode"},
	}

	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			if e.got != e.want {
				t.Errorf("Expected %s, got: %s", e.want, e.got)
			}
		})
	}
}

type invalidReader struct{}

func (i *invalidReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Some error")
}

func TestInvalidContentParse(t *testing.T) {
	table := []struct {
		name    string
		content io.Reader
		want    error
	}{
		{"invalid content", strings.NewReader("invalid content"), ErrCouldNotParseContent},
		{"nil content", nil, ErrCouldNotGetContent},
		{"invalid reader", &invalidReader{}, ErrCouldNotGetContent},
	}

	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			_, err := Parse(e.content)
			if err != e.want {
				t.Errorf("Expected %#v, but got: %#v", e.want, err)
			}
		})
	}
}

func TestSort(t *testing.T) {
	feed1, err := Parse(strings.NewReader(podcastfeed))
	if err != nil {
		t.Errorf("Didn't expect error but got: %#v", err)
	}
	feed2, err := Parse(strings.NewReader(podcastfeedReversed))
	if err != nil {
		t.Errorf("Didn't expect error but got: %#v", err)
	}

	if feed1.Channel.Items[0].Title.Title != feed2.Channel.Items[0].Title.Title {
		t.Errorf("Expected title to be the same after ordering, but it wasn't")
	}

	expectedDate, _ := time.Parse(Layout, "Thu, 21 Dec 2016 16:01:07 +0000")
	feed1Time, _ := time.Parse(Layout, feed1.Channel.Items[0].Date.Date)
	feed2Time, _ := time.Parse(Layout, feed2.Channel.Items[0].Date.Date)

	table := []struct {
		name string
		got  time.Time
		want time.Time
	}{
		{"time should be early to later for feed1", feed1Time, expectedDate},
		{"time should be early to later for feed2", feed2Time, expectedDate},
	}
	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			if e.got.Sub(e.want) > 1*time.Millisecond {
				t.Errorf("Expected ordering from early to later, got %#v, want: %#v", e.got.Format(Layout), e.want.Format(Layout))
			}
		})
	}
}
