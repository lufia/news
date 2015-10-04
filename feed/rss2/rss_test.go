package rss2

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
		for _, p := range feed.Channel.Items {
			t.Logf("\t%v\n", *p)
		}
	}
}

var xmlStringSimple = strings.TrimSpace(`
<?xml version='1.0' encoding='UTF-8'?>
<rss version='2.0'>
	<channel>
		<title>PHP &amp; JavaScript：更新情報</title>
		<link>http://phpjavascriptroom.com/</link>
		<description>PHP &amp; JavaScript Room：新着3件</description>
		<item>
			<title>記事タイトル3</title>
			<link>http://phpjavascriptroom.com/post3.html</link>
			<description>記事の内容です。</description>
			<pubDate>Wed, 11 Jun 2008 15:30:59 +0900</pubDate>
		</item>
		<item>
			<title>記事タイトル2</title>
			<link>http://phpjavascriptroom.com/post2.html</link>
			<description>記事の内容です。</description>
			<pubDate>Tue, 10 Jun 2008 15:30:59 +0900</pubDate>
		</item>
		<item>
			<title>記事タイトル1</title>
			<link>http://phpjavascriptroom.com/post1.html</link>
			<description>記事の内容です。</description>
			<pubDate>Mon, 9 Jun 2008 20:50:30 +0900</pubDate>
		</item>
	</channel>
</rss>
`)
