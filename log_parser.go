package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"time"
)

// LogParser is a ClientRequestSource which reads a proxy access log file and
// generates a new ClientRequest for each line to be replayed when stress
// testing a server.
type LogParser struct {
	// BufferSize determines how many requests may be queued before a client
	// worker is available to process the queue.
	BufferSize int

	// Replay determines how many times a log file will be replayed.
	// The default is 0 for a single read through.
	// To loop through a log file infinitely, set Replay to < 0.
	Replay int

	r *bufio.Reader
	l LogLexer
}

func NewLogParser(r io.Reader, l LogLexer) *LogParser {
	buf := bufio.NewReader(r)
	return &LogParser{r: buf, l: l, BufferSize: 4096}
}

func (c *LogParser) Get() (<-chan *ClientRequest, error) {
	if c.r == nil {
		return nil, fmt.Errorf("no reader defined for log parser")
	}

	if c.l == nil {
		return nil, fmt.Errorf("no lexer defined for log parser")
	}

	// create return channel
	ch := make(chan *ClientRequest, c.BufferSize)

	// go to work
	go func() {
		var err error
		var b []byte
		var line string
		var lineno uint64
		var isPrefix bool

		for {
			line = ""
			lineno++

			// create request
			r := &ClientRequest{
				Sequence: lineno,
			}

			// read in a full length line
			for isPrefix = true; isPrefix && err == nil; {
				b, isPrefix, err = c.r.ReadLine()
				if err == nil {
					if len(b) == 0 {
						// blank line... try again
						isPrefix = true
					} else {
						line += string(b)
					}
				}
			}

			// check for io errors or EOF
			if err != nil {
				if err == io.EOF {
					break
				} else {
					r.err = fmt.Errorf("line %d: %v", lineno, err)
					ch <- r
					continue
				}
			}

			// parse line
			m, err := c.l.Lex(line)
			if err != nil {
				r.err = fmt.Errorf("line %d: %v", lineno, err)
				ch <- r
				continue
			}

			// extract http method
			if method, ok := m["method"]; ok && method != "" {
				r.Method = method

				// TODO: send HTTPs log entries
				if method == "CONNECT" {
					continue
				}

			} else {
				r.err = fmt.Errorf("http method not found in input line %d", lineno)
				ch <- r
				continue
			}

			// extract url
			if u, ok := m["url"]; ok && u != "" {
				r.URL = u
			} else {
				r.err = fmt.Errorf("url not found in input line %d", lineno)
				ch <- r
				continue
			}

			// extract bytes
			if b, ok := m["bytes"]; ok && b != "" {
				bi, err := strconv.Atoi(b)
				if err != nil {
					r.err = fmt.Errorf("invalid byte count '%s' in input line %d", b, lineno)
					ch <- r
					continue
				}

				r.Bytes = uint64(bi)
			}

			// extract duration
			if d, ok := m["duration"]; ok && d != "" {
				di, err := strconv.Atoi(d)
				if err != nil {
					r.err = fmt.Errorf("invalid duration %dms in input line %d", d, lineno)
					ch <- r
					continue
				}

				r.Duration = time.Millisecond * time.Duration(di)
			}

			// ship it
			ch <- r
		}

		close(ch)
	}()

	return ch, nil
}
