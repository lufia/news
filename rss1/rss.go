package rss1

import (
	"encoding/xml"
	"io"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"RDF"`
	Channel *Channel `xml:"channel"`
	Items   []*Item  `xml:"item"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Date        time.Time `xml:"date"`
	Language    string    `xml:"language"`
	Indexes     []*Index  `xml:"items>Seq>li"`
}

type Index struct {
	URL string `xml:"resource,attr"`
}

type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Creator     string    `xml:"creator"`
	Date        time.Time `xml:"date"`
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
