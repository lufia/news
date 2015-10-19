package main

import (
	"log"
	"os"

	"github.com/lufia/news/feed"
)

func main() {
	f, err := feed.Parse(os.Stdin)
	if err != nil {
		log.Fatalln("Parse:", err)
	}
	log.Printf("%#v", *f)
}
