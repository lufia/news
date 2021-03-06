package news

import (
	"strings"
	"testing"
)

func TestDetectDialect(t *testing.T) {
	tab := []struct {
		xml  string
		want *Dialect
	}{
		{
			xml: `<?xml version="1.0"?>
				<feed xmlns="http://www.w3.org/2005/Atom">
				</feed>`,
			want: atomDialect,
		},
		{
			xml: `<?xml version="1.0"?>
				<feed version="0.3" xmlns="http://purl.org/atom/ns#" xmlns:dc="http://purl.org/dc/elements/1.1/">
				</feed>`,
			want: atomDialect,
		},
		{
			xml: `<?xml version="1.0"?>
				<rdf:RDF xmlns="http://purl.org/rss/1.0/"
					xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
					xmlns:dc="http://purl.org/dc/elements/1.1/"
					xmlns:content="http://purl.org/rss/1.0/modules/content/"
					xml:lang="ja">
				</rdf:RDF>`,
			want: rss1Dialect,
		},
		{
			xml: `<?xml version="1.0"?>
				<rss version="2.0">
				</rss>`,
			want: rss2Dialect,
		},
	}
	for _, v := range tab {
		r := strings.NewReader(v.xml)
		d, err := DetectDialect(r)
		if err != nil {
			t.Errorf("DetectDialect(%q) = %v", v.xml, err)
			continue
		}
		if d != v.want {
			t.Errorf("DetectDialect(%q) = %v; want %v", v.xml, d, v.want)
		}
	}
}

func TestCleanup(t *testing.T) {
	tab := []struct {
		s    string
		want string
	}{
		{s: "abc", want: "abc"},
		{s: "abc\v", want: "abc"},
		{s: "\vabc", want: "abc"},
		{s: "\va\vb\vc\v", want: "abc"},
		{s: "\v\v", want: ""},
	}
	for _, v := range tab {
		r := Cleanup([]byte(v.s))
		s := string(r)
		if s != v.want {
			t.Errorf("Cleanup(%q) = %q; want %q", v.s, s, v.want)
		}
	}
}
