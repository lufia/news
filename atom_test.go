package main

import (
	"strings"
	"testing"

	"github.com/lufia/news/feed/atom"
)

func TestNewArrival(t *testing.T) {
	r := strings.NewReader(xmlStringSimple)
	feed, err := atom.Parse(r)
	if err != nil {
		t.Fatalf("Parse(%q) = %v", xmlStringSimple, err)
	}

	tab := []struct {
		Time  time.Time
		Count int
	}{
		{Time: time.Date(2003, 12, 13, 18, 30, 03, 0, time.UTC), Count: 0},
		{Time: time.Date(2003, 12, 13, 18, 30, 02, 0, time.UTC), Count: 0},
		{Time: time.Date(2003, 12, 13, 18, 30, 01, 0, time.UTC), Count: 1},
	}
	for _, v := range tab {
		a := feed.NewArrival(v.Time)
		if len(a) != v.Count {
			t.Errorf("EntriesAfter(%v) = %d; Expect %d", v.Time, len(a), v.Count)
		}
	}
}
