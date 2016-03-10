package main

import (
	"errors"
	"fmt"
	"github.com/pivotal-golang/bytefmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// error to return for redirect requests
var errNoRedirect = errors.New("ignore redirects")

// Client may be used to stress test a server.
type Client struct {
	// Workers specifies the number of threads and HTTP connections to use for
	// stress testing.
	Workers int

	// StaggeryBy specifies how long to delay between starting each worker.
	StaggerBy time.Duration

	// KeepAlive specifies that the client should reuse TCP connections.
	KeepAlive bool

	proxyURL *url.URL
	source   ClientRequestSource
	sizzlin  bool
}

// NewClient returns the default configuration of a stress test client run.
func NewClient(proxyURL *url.URL, src ClientRequestSource) *Client {
	return &Client{
		proxyURL:  proxyURL,
		source:    src,
		KeepAlive: true,
		Workers:   8192,
		StaggerBy: time.Millisecond,
	}
}

// Sizzle will request lots of sausages from a sausage server using the Client
// configuration.
func (c *Client) Sizzle() error {
	if c.sizzlin {
		return fmt.Errorf("Already sizzlin'")
	}

	c.sizzlin = true

	// enforce min workers
	if c.Workers < 1 {
		c.Workers = 1
	}

	// create transport
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   c.Workers,
	}

	// configure http client
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errNoRedirect
		},
	}

	// parse proxy url
	if c.proxyURL != nil {
		tr.Proxy = http.ProxyURL(c.proxyURL)
	}

	// start the request producer
	reqs, err := c.source.Get()
	if err != nil {
		return err
	}

	// start stats collector
	resps := make(chan *ClientResponse, c.Workers*10)
	go c.stat(resps)

	// do requests
	done := make(chan bool, 0)
	for i := 0; i < c.Workers; i++ {
		go c.do(client, reqs, resps, done)

		if c.StaggerBy > 0 {
			time.Sleep(c.StaggerBy)
		}
	}

	// wait for workers to finish
	for i := 0; i < c.Workers; i++ {
		<-done
	}

	close(done)
	close(resps)

	c.sizzlin = false

	return nil
}

func (c *Client) stat(resps <-chan *ClientResponse) {
	var count, countPrev uint64
	var errs, errPrev uint64
	var data, dataPrev uint64

	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C: // timer tick - update stats
			// requests/sec
			rps := float64(count - countPrev)
			countPrev = count

			// errors/sec
			eps := float64(errs - errPrev)
			errPrev = errs

			// bits/sec
			bps := float64(data-dataPrev) * 8
			dataPrev = data

			fmt.Printf("Requests: %d (%.1f/sec) Transferred: %sb (%s) Errors: %d (%.1f/sec) Workers: %d\n", count, rps, bytefmt.ByteSize(data), bytefmt.BPSSize(bps), errs, eps, c.Workers)

		case resp := <-resps: // new response received
			if resp == nil {
				break
			}

			// increment
			count++
			data += resp.Size()
			if resp.Err() != nil {
				fmt.Fprintf(os.Stderr, "  %v in sequence %d\n", resp.Err(), resp.req.Sequence)
				errs++
			}
		}
	}
}

func (c *Client) do(client *http.Client, reqs <-chan *ClientRequest, resps chan<- *ClientResponse, done chan<- bool) {
	var err error
	var hreq *http.Request
	var hresp *http.Response
	var start time.Time

	for req := range reqs {
		start = time.Now()
		resp := &ClientResponse{req: req}

		if req.Method == "CONNECT" {
			// TODO: implement HTTPS requests
			resp.err = fmt.Errorf("CONNECT not implemented")

		} else {
			// reset error
			err = nil

			// create http request
			hreq, err = req.HTTPRequest()
			if err != nil {
				resp.err = err
			} else {
				// configure keep-alives
				if c.KeepAlive {
					hreq.Header.Set("Connection", "keep-alive")
				} else {
					hreq.Header.Set("Connection", "close")
				}

				// send request
				hresp, err = client.Do(hreq)
				if err == nil {
					// read response body
					defer hresp.Body.Close()

					resp.resp = hresp
					if b, err := ioutil.ReadAll(hresp.Body); err != nil {
						resp.err = err
					} else {
						resp.size = uint64(len(b))
						if req.Bytes > 0 && req.Method != "HEAD" && len(b) != int(req.Bytes) {
							resp.err = fmt.Errorf("expected %d bytes; got: %d in sequence %d", req.Bytes, len(b), req.Sequence)
						}
					}
				} else if uerr, ok := err.(*url.Error); ok && uerr.Err == errNoRedirect {
					resp.resp = hresp
				} else {
					resp.err = err
				}
			}
		}

		resp.duration = time.Now().Sub(start)
		resps <- resp
	}

	done <- true
}
