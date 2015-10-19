package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"net/smtp"
	"os"

	"github.com/lufia/news/feed"
)

var (
	addr     = flag.String("a", "localhost:587", "smtp address")
	user     = flag.String("u", "", "smtp username")
	password = flag.String("p", "", "smtp password")
	from     = flag.String("f", "", "from address")
	to       = flag.String("t", "", "to address")
)

func main() {
	f, err := feed.Parse(os.Stdin)
	if err != nil {
		log.Fatalln("Parse:", err)
	}
	m := MailMsg{
		From:    *from,
		To:      []string{*to},
		Article: f.Articles[0],
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = m.WriteTo(w)
	if err != nil {
		log.Fatalln("WriteTo:", err)
	}
	err = w.Flush()
	if err != nil {
		log.Fatalln("Flush:", err)
	}

	host, _, err := net.SplitHostPort(*addr)
	if err != nil {
		log.Fatalln("SplitHostPort:", err)
	}
	auth := smtp.PlainAuth("", *user, *password, host)
	err = smtp.SendMail(addr, auth, *from, []string{*to}, buf.Bytes())
	if err != nil {
		log.Fatalln("SendMail:", err)
	}
}
