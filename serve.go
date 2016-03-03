package main

import (
	"fmt"
	"github.com/pivotal-golang/bytefmt"
	"net/http"
	"strconv"
	"time"
)

var sausages = []byte("sausages\n")

func serveSausages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// set status code
	status := 200
	if s := r.Header.Get("X-Sausage-Status"); s != "" {
		i, err := strconv.Atoi(s)
		panicon(err)

		status = i
		w.WriteHeader(status)
	}

	// set server
	w.Header().Set("Server", "sausage/0.0.1")

	// set keep-alive
	if r.Header.Get("Connection") == "keep-alive" {
		w.Header().Set("Connection", "keep-alive")
	}

	// set mimetype
	if ct := r.Header.Get("X-Sausage-Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}

	// set content length
	contentLength := 9
	if s := r.Header.Get("X-Sausage-Content-Length"); s != "" {
		i, err := strconv.Atoi(s)
		panicon(err)

		if i > 0 {
			contentLength = i
		}
	}

	// set duration
	duration := 0
	if s := r.Header.Get("X-Sausage-Duration"); s != "" {
		i, err := strconv.Atoi(s)
		panicon(err)
		duration = i
	}

	// sleep
	if duration > 0 {
		time.Sleep(time.Millisecond * time.Duration(duration))
	}

	// render body
	n := 0
	i := 0
	var err error
	for written := 0; written < contentLength; written += n {
		i = len(sausages)
		if contentLength-written < i {
			i = contentLength - written
		}

		n, err = w.Write(sausages[:i])
		panicon(err)
	}

	// log
	now := time.Now()
	fmt.Printf("%v %s/%d %sb in %v\n", now, r.Method, status, bytefmt.ByteSize(uint64(contentLength)), now.Sub(start))
}
