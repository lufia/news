package rss2

import (
	"encoding/xml"
	"io"
	"time"
)

type Date time.Time

const (
	RFC2822 = "Mon, _2 Jan 2006 15:04:05 -0700"
)

func (date Date) String() string {
	return time.Time(date).String()
}

func (date *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var s string
	err = d.DecodeElement(&s, &start)
	if err != nil {
		return
	}
	t, err := time.Parse(RFC2822, s)
	if err != nil {
		return
	}
	*date = Date(t)
	return
}

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"` // dc:language
	Creator     string    `xml:"creator"`  // dc:creator
	Date        time.Time `xml:"date"`     // dc:date
	Items       []*Item   `xml:"item"`
}

type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Category    string    `xml:"category,omitempty"`
	Subject     string    `xml:"subject"` // dc:subject
	Creator     string    `xml:"creator"` // dc:creator
	PubDate     Date      `xml:"pubDate"`
	Date        time.Time `xml:"date"` // dc:date
}

func Parse(r io.Reader) (feed *Feed, err error) {
	var x Feed
	d := xml.NewDecoder(r)
	err = d.Decode(&x)
	if err != nil {
		return
	}
	feed = &x
	return
}
