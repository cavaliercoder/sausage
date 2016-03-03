package main

import (
	"strings"
	"testing"
)

func TestLogParser(t *testing.T) {
	fields := []string{"", "", "size", "", "", "duration", "method", "url", "", "", "mime_type", "agent"}
	p := NewLogParser(strings.NewReader(ssvLog), NewSSVLexer(fields))

	ch, err := p.Get()
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	for r := range ch {
		if r.Err() != nil {
			t.Errorf("%v", r.Err())
		}
	}
}
