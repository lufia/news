package main

import (
	"html/template"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/lufia/news/feed"
)

type MailMsg struct {
	Article feed.Article
	From    string
	To      []string
}

type errWriter struct {
	err error
	w   io.Writer
}

func (w *errWriter) Write(p []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write(p)
}

func (w *errWriter) Println(args ...string) {
	for i, s := range args {
		if i > 0 {
			w.Write([]byte{' '})
		}
		w.Write([]byte(s))
	}
	w.Write([]byte{'\n'})
}

func (msg *MailMsg) WriteTo(w io.Writer) (err error) {
	m := multipart.NewWriter(w)

	fout := errWriter{w: w}
	fout.Println("Subject:", msg.Article.Title)
	fout.Println("From:", msg.From)
	fout.Println("To:", strings.Join(msg.To, ", "))
	ctype := "multipart/alternative; boundary=" + m.Boundary()
	fout.Println("Content-Type:", ctype)
	fout.Println()

	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", "text/html; charset=UTF-8")
	w1, err := m.CreatePart(h)
	if err != nil {
		return
	}
	body := template.HTML(msg.Article.Content)
	err = msgTemplate.Execute(w1, body)
	return
}

var msgContainer = strings.TrimSpace(`
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
</head>
<body>

<table width="580" cellpadding="0" cellspacing="0" border="0" align="center">
<tr>
<td>
{{.}}
</td>
</tr>
</table>

</body>
</html>
`)

var (
	msgTemplate *template.Template
)

func init() {
	t := template.New("mail")
	msgTemplate = template.Must(t.Parse(msgContainer))
}
