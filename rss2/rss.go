package rss2

import (
	"encoding/xml"
	"errors"
	"io"
	"time"
)

type Date time.Time

const (
	RFC2822  = "Mon, _2 Jan 2006 15:04:05 -0700"
	RFC2822Z = "Mon, _2 Jan 2006 15:04:05 MST"
)

var (
	errNoItemID = errors.New("item hasn't <guid> or <link> tag")
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
		if t, err = time.Parse(RFC2822Z, s); err != nil {
			return
		}
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
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language,omitempty"`
	LastBuildDate Date     `xml:"lastBuildDate,omitempty"`
	Category      Category `xml:"category,omitempty"`
	Items         []*Item  `xml:"item"`

	Creator string    `xml:"creator"` // dc:creator
	Date    time.Time `xml:"date"`    // dc:date
}

type Item struct {
	Title       string     `xml:"title,omitempty"`
	Link        string     `xml:"link,omitempty"`
	Description string     `xml:"description,omitempty"`
	Author      string     `xml:"author,omitempty"` // author's email address
	Categories  []Category `xml:"category,omitempty"`
	Guid        Guid       `xml:"guid,omietmpty"`
	PubDate     Date       `xml:"pubDate,omitempty"`

	Subject string    `xml:"subject,omitempty"` // dc:subject
	Creator string    `xml:"creator,omitempty"` // dc:creator
	Date    time.Time `xml:"date,omitempty"`    // dc:date
	Encoded string    `xml:"encoded,omitempty"` // content:encoded
}

func (item *Item) Content() string {
	if item.Encoded != "" {
		return item.Encoded
	}
	return item.Description
}

func (item *Item) ID() (string, error) {
	if item.Guid.IsPermaLink && item.Guid.Content != "" {
		return item.Guid.Content, nil
	}
	if item.Link != "" {
		return item.Link, nil
	}
	return "", errNoItemID
}

type Category struct {
	Domain  string `xml:"domain,attr,omitempty"`
	Content string `xml:",chardata"`
}

type Guid struct {
	IsPermaLink bool   `xml:"isPermaLink,omitempty"`
	Content     string `xml:",chardata"`
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
