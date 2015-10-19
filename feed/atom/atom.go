package atom

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
	"time"

	"golang.org/x/net/html"
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

// IsZeroはtが空だった場合にtrueを返す。
func (t Text) IsZero() bool {
	return t.Content == ""
}

func buildPlain(n *html.Node) (s string, err error) {
	buf := new(bytes.Buffer)
	err = html.Render(buf, n)
	if err != nil {
		return
	}
	s = buf.String()
	return
}

func (t Text) Plain() (s string, err error) {
	switch t.Type {
	case "html", "xhtml":
		err = errors.New("not implement")
	case "text":
		s = t.Content
	default:
		s = t.Content
	}
	return
}

func (t Text) HTML() (s string, err error) {
	switch t.Type {
	case "html":
		t := html.UnescapeString(t.Content)
		s = fmt.Sprintf("<div>%s</div>", t)
	case "xhtml":
		r := strings.NewReader(t.Content)
		tokenizer := html.NewTokenizer(r)
		err = nextToken(tokenizer)
		if err != nil {
			return
		}
		s, err = buildHTML(tokenizer)
	case "text":
		s = fmt.Sprintf("<pre>%s</pre>", t.Content)
	default:
		s = fmt.Sprintf("<pre>%s</pre>", t.Content)
	}
	return
}

func nextToken(tokenizer *html.Tokenizer) error {
	if t := tokenizer.Next(); t == html.ErrorToken {
		return tokenizer.Err()
	}
	return nil
}

func buildHTML(tokenizer *html.Tokenizer) (s string, err error) {
	buf := new(bytes.Buffer)

	bp := 0
	if tag, _ := tokenizer.TagName(); string(tag) == "div" {
		div := tokenizer.Raw()
		buf.Write(div)
		bp = len(div)
		err = nextToken(tokenizer)
	}

	ep := bp
	for err != io.EOF {
		if err != nil && err != io.EOF {
			return
		}
		ep = buf.Len()
		b := tokenizer.Raw()
		if _, err := buf.Write(b); err != nil {
			return "", err
		}
		err = nextToken(tokenizer)
	}
	b := buf.Bytes()
	if bp > 0 {
		b = b[bp:ep]
	}
	return string(b), nil
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

	//Contributors []Person `xml:"contributor,omitempty"`
	//Generator Generator `xml:"generator,omitempty"`
	//Icon string `xml:"icon,omitempty"`
	//Logo string `xml:"logo,omitempty"`

	Title      Text       `xml:"title"`
	Subtitle   Text       `xml:"subtitle,omitempty"`
	Links      []Link     `xml:"link"`
	Authors    []Person   `xml:"author"`
	ID         string     `xml:"id"`
	Rights     Text       `xml:"rights,omitempty"`
	Updated    time.Time  `xml:"updated"`
	Summary    string     `xml:"summary,omitempty"`
	Categories []Category `xml:"category,omitempty"`
	Entries    []*Entry   `xml:"entry"`
}

// EntryはAtom文書におけるEntry要素をあらわす。
type Entry struct {
	//Contributors []Person `xml:"contributor,omitempty"`
	//Created time.Time `xml:"created,omitempty"?
	//Source Link?

	Title      Text       `xml:"title"`
	Links      []Link     `xml:"link,omitempty"`
	Authors    []Person   `xml:"author,omitempty"`
	Categories []Category `xml:"category,omitempty"`
	ID         string     `xml:"id"`
	Updated    time.Time  `xml:"updated"`
	Published  time.Time  `xml:"published,omitempty"`
	Rights     Text       `xml:"rights,omitempty"`
	Summary    Text       `xml:"summary,omitempty"`
	Content    Text       `xml:"content,omitempty"`
}

func (entry *Entry) Article() string {
	return ""
}

type MailBody Entry

func (body *MailBody) WriteTo(w io.Writer) (n int64, err error) {
	m := multipart.NewWriter(w)
	err = m.SetBoundary("multipart_boundary_str")
	if err != nil {
		return
	}
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
