package feed

import (
	"encoding/xml"
	"errors"
	"io"
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

type Dialect string

var decisionTable = []struct {
	elem    distinctElement
	dialect Dialect
}{
	{
		elem: distinctElement{
			XMLName: xml.Name{
				Space: "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
				Local: "RDF",
			},
		},
		dialect: Dialect("rss1.0"),
	},
	{
		elem: distinctElement{
			XMLName: xml.Name{
				Local: "rss",
			},
			Version: "2.0",
		},
		dialect: Dialect("rss2.0"),
	},
	{
		elem: distinctElement{
			XMLName: xml.Name{
				Space: "http://www.w3.org/2005/Atom",
				Local: "feed",
			},
		},
		dialect: Dialect("atom"),
	},
}

var (
	errUnknownDialect = errors.New("unknown dialect")
)

func DetectDialect(r io.Reader) (Dialect, error) {
	var x distinctElement
	d := xml.NewDecoder(r)
	if err := d.Decode(&x); err != nil {
		return "", err
	}
	for _, v := range decisionTable {
		if x == v.elem {
			return v.dialect, nil
		}
	}
	return "", errUnknownDialect
}
