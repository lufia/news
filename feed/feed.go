// Package feed presets rss/atom reader.
package feed

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"

	"github.com/lufia/news/feed/atom"
	"github.com/lufia/news/feed/rss1"
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
			return nil, errors.New("not implemented")
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

func Parse(r io.Reader) (feed interface{}, err error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
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
