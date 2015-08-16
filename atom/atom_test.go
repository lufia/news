package atom

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tab := []struct {
		XMLString    string
		ExpectedFeed Feed
	}{
		{
			XMLString: xmlStringSimple,
			ExpectedFeed: Feed{
				XMLName: xml.Name{
					Space: "http://www.w3.org/2005/Atom",
					Local: "feed",
				},
				Title:   S("Example Feed"),
				Links:   URLs("http://example.org/"),
				Updated: time.Date(2003, 12, 13, 18, 30, 02, 0, time.UTC),
				Authors: Persons("John Doe"),
				ID:      "urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6",
				Entries: []*Entry{
					&Entry{
						Title:   S("Atom-Powered Robots Run Amok"),
						Links:   URLs("http://example.org/2003/12/13/atom03"),
						ID:      "urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a",
						Updated: time.Date(2003, 12, 13, 18, 30, 02, 0, time.UTC),
						Summary: S("Some text."),
					},
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
		expect := &v.ExpectedFeed
		if !reflect.DeepEqual(feed, expect) {
			t.Errorf("Parse(%q) = %#v; Expect %#v", v.XMLString, feed, expect)
		}
	}
}

func TestText_Plain(t *testing.T) {
	tab := []struct {
		Text   Text
		Expect string
	}{
		{
			Text:   Text{Type: "text", Content: "<test>ab</test>"},
			Expect: "<test>ab</test>",
		},
		//{
		//	Text:   Text{Type: "html", Content: "&lt;p&gt;ab&lt;/p&gt;"},
		//	Expect: "ab",
		//},
		//{
		//	Text:   Text{Type: "xhtml", Content: "<p>&lt;ab&gt;</p>"},
		//	Expect: "<ab>",
		//},
	}
	for _, v := range tab {
		s, err := v.Text.Plain()
		if err != nil {
			t.Fatalf("(%#v).Plain() = %v", v.Text, err)
		}
		if s != v.Expect {
			t.Errorf("(%#v).Plain() = %q; Expect %q", v.Text, s, v.Expect)
		}
	}
}

func TestText_HTML(t *testing.T) {
	tab := []struct {
		Text   Text
		Expect string
	}{
		{
			Text:   Text{Type: "text", Content: "<test>ab</test>"},
			Expect: "<pre><test>ab</test></pre>",
		},
		{
			Text:   Text{Type: "html", Content: "&lt;test&gt;ab&lt;/test&gt;"},
			Expect: "<div><test>ab</test></div>",
		},
		{
			Text:   Text{Type: "xhtml", Content: "<div>&lt;em&gt;ab</div>"},
			Expect: "&lt;em&gt;ab",
		},
		{
			Text:   Text{Type: "xhtml", Content: "<div><em>ab</em></div>"},
			Expect: "<em>ab</em>",
		},
	}
	for _, v := range tab {
		s, err := v.Text.HTML()
		if err != nil {
			t.Fatalf("(%#v).HTML() = %v", v.Text, err)
		}
		if s != v.Expect {
			t.Errorf("(%#v).HTML() = %q; Expect %q", v.Text, s, v.Expect)
		}
	}
}

func TestMailBody_WriteTo(t *testing.T) {
	tab := []struct {
		Entry          *Entry
		ExpectedString string
	}{
		{
			Entry: &Entry{
				Content: S("test"),
			},
		},
	}
	for _, v := range tab {
		var buf bytes.Buffer
		body := (*MailBody)(v.Entry)
		_, err := body.WriteTo(&buf)
		if err != nil {
			t.Fatalf("WriteTo() = %v", err)
		}
		t.Log("Output:", string(buf.Bytes()))
	}
}

var xmlStringSimple = strings.TrimSpace(`
<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
	<title>Example Feed</title>
	<link href="http://example.org/"/>
	<updated>2003-12-13T18:30:02Z</updated>
	<author>
		<name>John Doe</name>
	</author>
	<id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>
	<entry>
		<title>Atom-Powered Robots Run Amok</title>
		<link href="http://example.org/2003/12/13/atom03"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
	</entry>
</feed>
`)
