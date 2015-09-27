package main

import (
	"io"

	"github.com/lufia/news/feed/atom"
)

func (h *atomMailHeader) WriteTo(w io.Writer) (n int64, err error) {
}

func (body *atomMailBody) WriteTo(w io.Writer) (n int64, err error) {
}
