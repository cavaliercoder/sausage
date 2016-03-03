package main

import (
	"fmt"
	"regexp"
	"strings"
)

// SSVLexer is a lexer for space-separated values.
type SSVLexer struct {
	fields []string
}

// regex pattern to split a space-separated values row into an array
var ssvPattern = regexp.MustCompile(`("(?:[^"\\]|\\.)*")|([^\s]+)`)

func NewSSVLexer(fields []string) LogLexer {
	return &SSVLexer{fields: fields}
}

func (c *SSVLexer) Lex(s string) (map[string]string, error) {
	// split line into array
	cols := ssvPattern.FindAllString(s, -1)
	if len(cols) < len(c.fields) {
		return nil, fmt.Errorf("expected at least %d fields; got %d in: %s", len(c.fields), len(cols), s)
	}

	// TODO: add modulators to SSVLexer to allow for fields such as TCP_MISS/200

	// map each column
	m := make(map[string]string, len(c.fields))
	for i := 0; i < len(c.fields); i++ {
		if c.fields[i] != "" {
			// assign and strip quotes
			m[c.fields[i]] = strings.Trim(cols[i], "\"")
		}
	}

	return m, nil
}
