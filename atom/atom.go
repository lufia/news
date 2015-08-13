package atom

import (
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/textproto"
	"time"
)

// CategoryはAtom文書におけるCategory要素をあらわす。
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

func (t Text) IsZero() bool {
	return t.Content == ""
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

// FeedはAtom文書におけるFeed要素をあらわす。
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
	Entries  []*Entry  `xml:"entry"`
}

// EntryはAtom文書におけるEntry要素をあらわす。
type Entry struct {
	//Contributors []Person `xml:"contributor,omitempty"`
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

type MailBody Entry

func (body *MailBody) WriteTo(w io.Writer) (n int64, err error) {
	m := multipart.NewWriter(w)
	written, err := body.writeTextTo(m)
	if err != nil {
		return
	}
	n += written
	written, err = body.writeHTMLTo(m)
	if err != nil {
		return
	}
	n += written
	return
}

func (body *MailBody) textBody() []byte {
	if !body.Content.IsZero() {
		return []byte(body.Content.Content)
	}
	return []byte(body.Summary.Content)
}

func (body *MailBody) htmlBody() []byte {
	return body.textBody() // TODO: quick
}

func (body *MailBody) writeTextTo(m *multipart.Writer) (n int64, err error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", "text/plain")
	w, err := m.CreatePart(h)
	if err != nil {
		return
	}
	written, err := w.Write(body.textBody())
	if err != nil {
		return
	}
	n += int64(written)
	return
}

func (body *MailBody) writeHTMLTo(m *multipart.Writer) (n int64, err error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", "text/html")
	w, err := m.CreatePart(h)
	if err != nil {
		return
	}
	written, err := w.Write(body.htmlBody())
	if err != nil {
		return
	}
	n += int64(written)
	return
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
