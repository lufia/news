// Package feed presets rss/atom reader.
package news

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"time"
	"unicode/utf8"

	"github.com/lufia/news/atom"
	"github.com/lufia/news/rss1"
	"github.com/lufia/news/rss2"
)

type distinctElement struct {
	XMLName xml.Name
	Version string `xml:"version,attr"`
}

func (rule distinctElement) Match(v distinctElement) bool {
	x1 := rule.XMLName
	x2 := v.XMLName
	if x1.Space != "" && x1.Space != x2.Space {
		return false
	}
	if x1.Local != x2.Local {
		return false
	}
	if rule.Version != "" && rule.Version != v.Version {
		return false
	}
	return true
}

type Dialect struct {
	Type  string
	Parse func(r io.Reader) (feed interface{}, err error)
}

var (
	rss1Dialect = &Dialect{
		Type: "rss1.0",
		Parse: func(r io.Reader) (feed interface{}, err error) {
			return rss1.Parse(r)
		},
	}
	rss2Dialect = &Dialect{
		Type: "rss2.0",
		Parse: func(r io.Reader) (feed interface{}, err error) {
			return rss2.Parse(r)
		},
	}
	atomDialect = &Dialect{
		Type: "atom",
		Parse: func(r io.Reader) (feed interface{}, err error) {
			return atom.Parse(r)
		},
	}
)

var decisionTable = []struct {
	elem    distinctElement
	dialect *Dialect
}{
	{
		elem: distinctElement{
			XMLName: xml.Name{
				Space: "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
				Local: "RDF",
			},
		},
		dialect: rss1Dialect,
	},
	{
		elem: distinctElement{
			XMLName: xml.Name{
				Local: "rss",
			},
			Version: "2.0",
		},
		dialect: rss2Dialect,
	},
	{
		elem: distinctElement{
			XMLName: xml.Name{
				Space: "http://www.w3.org/2005/Atom",
				Local: "feed",
			},
		},
		dialect: atomDialect,
	},
	{
		elem: distinctElement{
			XMLName: xml.Name{
				Space: "http://purl.org/atom/ns#",
				Local: "feed",
			},
		},
		dialect: atomDialect,
	},
}

var (
	errUnknownDialect = errors.New("unknown dialect")
)

func DetectDialect(r io.Reader) (*Dialect, error) {
	var x distinctElement
	d := xml.NewDecoder(r)
	if err := d.Decode(&x); err != nil {
		return nil, err
	}
	for _, v := range decisionTable {
		if v.elem.Match(x) {
			return v.dialect, nil
		}
	}
	return nil, errUnknownDialect
}

// Cleanup discards invalid chars in XML 1.0.
func Cleanup(p []byte) []byte {
	tab := []rune{'\v'}
	for _, c := range tab {
		n := utf8.RuneLen(c)
		for {
			i := bytes.IndexRune(p, c)
			if i < 0 {
				break
			}
			copy(p[i:], p[i+n:])
			p = p[:len(p)-n]
		}
	}
	return p
}

func parse(r io.Reader) (feed interface{}, err error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	buf = Cleanup(buf)
	fin := bytes.NewReader(buf)
	d, err := DetectDialect(fin)
	if err != nil {
		return
	}
	_, err = fin.Seek(0, 0)
	if err != nil {
		return
	}
	return d.Parse(fin)
}

type Feed struct {
	Title    string
	URL      string
	Summary  string
	Articles []*Article
}

type Article struct {
	Title      string
	ID         string
	URL        string
	Authors    []string
	Published  time.Time
	Categories []string
	Content    string
}

func Parse(r io.Reader) (feed *Feed, err error) {
	p, err := parse(r)
	if err != nil {
		return
	}
	feed = &Feed{}
	switch v := p.(type) {
	case *rss1.Feed:
		err = feed.ImportFromRSS1(v)
		return
	case *rss2.Feed:
		err = feed.ImportFromRSS2(v)
		return
	case *atom.Feed:
		err = feed.ImportFromAtom(v)
		return
	default:
		return nil, errors.New("unknown feed type")
	}
}

func (feed *Feed) ImportFromRSS1(r *rss1.Feed) (err error) {
	feed.Title = r.Channel.Title
	feed.URL = r.Channel.Link
	feed.Summary = r.Channel.Description
	feed.Articles = make([]*Article, len(r.Items))
	for i, item := range r.Items {
		p := &Article{
			Title:     item.Title,
			ID:        r.Channel.Indexes[i].URL,
			URL:       item.Link,
			Authors:   []string{item.Creator},
			Published: item.Date,
			Content:   item.Description,
		}
		feed.Articles[i] = p
	}
	return nil
}

type rss2Item rss2.Item

func (v *rss2Item) Authors() []string {
	if v.Author == "" {
		return []string{}
	}
	return []string{v.Author}
}

func (v *rss2Item) Published() time.Time {
	return time.Time(v.PubDate)
}

func (feed *Feed) ImportFromRSS2(r *rss2.Feed) (err error) {
	feed.Title = r.Channel.Title
	feed.URL = r.Channel.Link
	feed.Summary = r.Channel.Description
	feed.Articles = make([]*Article, len(r.Channel.Items))
	for i, item := range r.Channel.Items {
		v := (*rss2Item)(item)
		p := &Article{
			Title:     item.Title,
			URL:       item.Link,
			Authors:   v.Authors(),
			Published: v.Published(),
			Content:   item.Content(),
		}
		if p.ID, err = item.ID(); err != nil {
			return
		}
		feed.Articles[i] = p
	}
	return
}

func (feed *Feed) ImportFromAtom(r *atom.Feed) (err error) {
	feed.Title = r.Title.Content
	feed.URL = r.AlternateURL()
	feed.Summary = r.Summary
	feed.Articles = make([]*Article, len(r.Entries))
	for i, entry := range r.Entries {
		p := &Article{
			Title:     entry.Title.Content,
			ID:        entry.ID,
			URL:       entry.AlternateURL(),
			Authors:   feed.atomAuthors(entry.Authors),
			Published: entry.PublishedTime(),
		}
		var s string
		s, err = entry.Content.HTML()
		if err != nil {
			return
		}
		p.Content = s
		feed.Articles[i] = p
	}
	return
}

func (feed *Feed) atomAuthors(authors []atom.Person) []string {
	a := make([]string, len(authors))
	for i, p := range authors {
		a[i] = p.Name
	}
	return a
}
