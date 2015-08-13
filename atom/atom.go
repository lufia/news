package atom

import (
	"encoding/xml"
	"io"
	"reflect"
	"time"
)

// CategoryはAtom文書におけるCategory要素
type Category struct {
	Term   string `xml:"term,attr"`
	Scheme string `xml:"scheme,attr,omitempty"`
	Label  string `xml:"label,attr,omitempty"`
}

// TextはAtom文書におけるTextコンストラクトをあらわす。
type Text struct {
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// PersonはAtom文書におけるPersonコンストラクトをあらわす。
type Person struct {
	Name  string `xml:"name"`
	URL   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
}

// LinkはAtom文書におけるLinkコンストラクトをあらわす。
type Link struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	URL  string `xml:"href,attr"`
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`

	//Version string `xml:"version,attr"`
	//Lang string `xml:"lang,attr,omitempty"`

	//Categories []Text `xml:"category,omitempty"`
	//Contributors []Person `xml:"contributor,omitempty"`
	//Generator Generator `xml:"generator,omitempty"`
	//Icon string `xml:"icon,omitempty"`
	//Logo string `xml:"logo,omitempty"`
	//Categories []Category `xml:"category,omitempty"`

	Title    Text      `xml:"title"`
	Subtitle Text      `xml:"subtitle,omitempty"`
	Links    []Link    `xml:"link"`
	Authors  []Person  `xml:"author"`
	ID       string    `xml:"id"`
	Rights   Text      `xml:"rights,omitempty"`
	Updated  time.Time `xml:"updated"`
	Summary  string    `xml:"summary,omitempty"`
	Entries  []Entry   `xml:"entry"`
}

type Entry struct {
	//Contributors []Person `xml:"contributor"`
	//Created time.Time `xml:"created,omitempty"?
	//Categories []Category `xml:"category,omitempty"`
	//Source Link?

	Title     Text      `xml:"title"`
	Links     []Link    `xml:"link,omitempty"`
	Author    []Person  `xml:"author,omitempty"`
	ID        string    `xml:"id"`
	Updated   time.Time `xml:"updated"`
	Published time.Time `xml:"published,omitempty"`
	Rights    Text      `xml:"rights,omitempty"`
	Summary   Text      `xml:"summary,omitempty"`
	Content   Text      `xml:"content,omitempty"`
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

func (feed *Feed) Equal(other *Feed) bool {
	return reflect.DeepEqual(feed, other)
}
