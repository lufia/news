package atom

import (
	"testing"
)

func TestFeedAlternateURL(t *testing.T) {
	tab := []struct {
		Links []Link
		URL   string
	}{
		{
			Links: []Link{
				{Rel: "alternate", Type: "text/html", URL: "http://localhost/"},
			},
			URL: "http://localhost/",
		},
		{
			Links: []Link{
				{Rel: "related", Type: "text/html", URL: "http://example.com/"},
				{Rel: "alternate", URL: "http://localhost/"},
			},
			URL: "http://localhost/",
		},
		{
			Links: []Link{
				{Type: "text/html", URL: "http://localhost/"},
			},
			URL: "http://localhost/",
		},
	}
	for _, v := range tab {
		var feed Feed
		feed.Links = v.Links
		s := feed.AlternateURL()
		if s != v.URL {
			t.Errorf("AlternateURL() = %q; want %q", s, v.URL)
		}
	}
}
