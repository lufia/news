package rss1

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tab := []struct {
		XMLString string
		Expected  *Feed
	}{
		{
			XMLString: xmlStringSimple,
			Expected: &Feed{
				XMLName: xml.Name{
					Space: "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
					Local: "RDF",
				},
			},
		},
	}
	for _, v := range tab {
		r := strings.NewReader(v.XMLString)
		feed, err := Parse(r)
		if err != nil {
			t.Errorf("Parse(%q) = %v", v.XMLString, err)
			continue
		}
		t.Logf("%v\n", *feed)
		t.Logf("%v\n", feed.Channel)
		for _, p := range feed.Items {
			t.Logf("\t%v\n", *p)
		}
	}
}

var xmlStringSimple = strings.TrimSpace(`
<?xml version="1.0" encoding="UTF-8"?>
<rdf:RDF xmlns="http://purl.org/rss/1.0/"
	xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xml:lang="ja">

	<channel rdf:about="サイトのRSSのURL">
		<title>サイトのタイトル</title>
		<link>サイトのURL</link>
		<description>サイトの内容</description>
		<dc:date>2015-02-01T02:03:45Z</dc:date>
		<dc:language>ja</dc:language> 
		<items>
		<rdf:Seq>
		<rdf:li rdf:resource="記事1のURL" />
		<rdf:li rdf:resource="記事2のURL" />
		</rdf:Seq>
		</items>
	</channel>

	<item rdf:about="記事1のURL">
		<title>記事1のタイトル</title>
		<link>記事1のURL</link>
		<description><![CDATA[記事1の内容]]></description>
		<dc:creator>記事1の作者名</dc:creator>
		<dc:date>2015-01-01T02:03:45+09:00</dc:date>
	</item>

	<item rdf:about="記事2のURL">
		<title>記事2のタイトル</title>
		<link>記事2のURL</link>
		<description><![CDATA[記事2の内容]]></description>
		<dc:creator>記事2の作者名</dc:creator>
		<dc:date>2015-02-01T02:03:45+09:00</dc:date>
	</item>
</rdf:RDF>
`)
