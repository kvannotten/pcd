package pcd

import (
	"net/http"
	"net/http/httptest"
)

var Podcastfeed = `<?xml version="1.0" encoding="UTF-8"?>
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
<!--END REPEAT--> 
   
</channel>
</rss>`

func testServer() *httptest.Server {
	return testServerWithBasicAuth("", "")
}

func testServerWithBasicAuth(username, password string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username != "" {
			u, p, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if u != username || p != password {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		if _, err := w.Write([]byte(Podcastfeed)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func testServerWithStatusCode(code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}))
}
