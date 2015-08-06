package atom

import (
	"errors"
	"io"
)

type Feed struct {
	title string
	url   string
}

func Parse(r io.Reader) (feed *Feed, err error) {
	err = errors.New("teest")
	return
}

func (feed *Feed) Title() string {
	return feed.title
}
