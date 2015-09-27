// Package feed presets rss/atom reader.
package feed

import (
	"encoding/xml"
	"errors"
	"io"
	"time"

	_ "github.com/lufia/news/feed/atom"
	_ "github.com/lufia/news/feed/rss1"
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
	Type string
}

var (
	rss1Dialect = &Dialect{
		Type: "rss1.0",
	}
	rss2Dialect = &Dialect{
		Type: "rss2.0",
	}
	atomDialect = &Dialect{
		Type: "atom",
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

type Info struct {
	Title   string
	URL     string
	Updated time.Time
}

type Channel interface {
	Articles() []Article
}

type Article interface {
	Info() *Info
	io.WriterTo
}
