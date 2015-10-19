package main

import (
	"testing"
)

func TestMailMsg_WriteTo(t *testing.T) {
	tab := []struct {
		Entry          *Entry
		ExpectedString string
	}{
		{
			Entry: &Entry{
				Content: S("test"),
			},
		},
	}
	for _, v := range tab {
		var buf bytes.Buffer
		msg := (*MailMsg)(v.Entry)
		_, err := msg.WriteTo(&buf)
		if err != nil {
			t.Fatalf("WriteTo() = %v", err)
		}
		t.Log("Output:", string(buf.Bytes()))
	}
}
