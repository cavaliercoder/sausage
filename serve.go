package main

import (
	"fmt"
	"github.com/pivotal-golang/bytefmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ServeSausages struct {
	server *http.Server
	stats  chan *ServerStat
	start  time.Time
}

var sausages = []byte("sausages\n")

func NewSausageServer() *ServeSausages {
	srv := &ServeSausages{
		stats: make(chan *ServerStat, 4096),
	}

	srv.server = &http.Server{
		Addr:        fmt.Sprintf("%s:%d", flagListenIface, flagListenPort),
		Handler:     srv,
		ConnState:   srv.ConnState,
		ReadTimeout: time.Second,
	}

	return srv
}

func (c *ServeSausages) ListenAndServe() error {
	c.start = time.Now()
	go c.stat()
	return c.server.ListenAndServe()
}

func (c *ServeSausages) stat() {
	t := time.NewTicker(time.Second)
	totals := &ServerStat{}
	last := &ServerStat{}
	for {
		select {
		case <-t.C:
			fmt.Printf("Requests: %d (%d/sec) Connections: %d (%d/sec) Transferred: %sb (%s)\n", totals.Requests, totals.Requests-last.Requests, totals.ConnectionsOpen, totals.ConnectionsOpen-last.ConnectionsOpen, bytefmt.ByteSize(uint64(totals.BytesTansferred)), bytefmt.BPSSize(float64(totals.BytesTansferred-last.BytesTansferred)))
			*last = *totals

		case st := <-c.stats:
			totals.Add(st)
		}
	}
}

func (c *ServeSausages) ConnState(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		c.stats <- &ServerStat{Connections: 1, ConnectionsOpen: 1}

	case http.StateHijacked, http.StateClosed:
		c.stats <- &ServerStat{ConnectionsOpen: -1}
	}
}

func (c *ServeSausages) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	} else {
		w.Header().Set("Connection", "close")
	}

	// set mimetype
	if ct := r.Header.Get("X-Sausage-Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}

	// sequence id
	seq := ""
	if seq = r.Header.Get("X-Sausage-Sequence"); seq != "" {
		w.Header().Set("X-Sausage-Sequence", seq)
	}

	// set content length
	contentLength := 9
	written := 0
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
	written = 0
	if r.Method != "HEAD" {
		n := 0
		i := 0
		var err error
		for ; written < contentLength; written += n {
			i = len(sausages)
			if contentLength-written < i {
				i = contentLength - written
			}

			n, err = w.Write(sausages[:i])
			panicon(err)
		}
	}

	// log
	now := time.Now()
	//fmt.Printf("%v %s/%d %sb in %v (sequence: %s)\n", now, r.Method, status, bytefmt.ByteSize(uint64(written)), now.Sub(start), seq)

	c.stats <- &ServerStat{BytesTansferred: int64(written), Requests: 1, TimeInRequest: now.Sub(start)}
}
