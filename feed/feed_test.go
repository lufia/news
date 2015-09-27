package feed

import (
	"strings"
	"testing"
)

func TestDetectDialect(t *testing.T) {
	tab := []struct {
		xml  string
		want Dialect
	}{
		{
			xml: `<?xml version="1.0"?>
				<feed xmlns="http://www.w3.org/2005/Atom">
				</feed>`,
			want: Dialect("atom"),
		},
		{
			xml: `<?xml version="1.0"?>
				<rdf:RDF xmlns="http://purl.org/rss/1.0/"
					xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
					xmlns:dc="http://purl.org/dc/elements/1.1/"
					xmlns:content="http://purl.org/rss/1.0/modules/content/"
					xml:lang="ja">
				</rdf:RDF>`,
			want: Dialect("rss1.0"),
		},
		{
			xml: `<?xml version="1.0"?>
				<rss version="2.0">
				</rss>`,
			want: Dialect("rss2.0"),
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
