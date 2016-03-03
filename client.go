package main

import (
	"fmt"
	"github.com/pivotal-golang/bytefmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Client may be used to stress test a server.
type Client struct {
	Workers int

	proxyURL *url.URL
	source   ClientRequestSource
	sizzlin  bool
}

// NewClient returns the default configuration of a stress test client run.
func NewClient(proxyURL *url.URL, src ClientRequestSource) *Client {
	return &Client{
		proxyURL: proxyURL,
		source:   src,
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
	resps := make(chan *ClientResponse, c.Workers*100)
	go c.stat(resps)

	// do requests
	done := make(chan bool, 0)
	for i := 0; i < c.Workers; i++ {
		go c.do(client, reqs, resps, done)
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

			fmt.Printf("Requests: %d (%.1f/sec) Transferred: %sb (%s) Errors: %d (%.1f/sec)\n", count, rps, bytefmt.ByteSize(data), bytefmt.BPSSize(bps), errs, eps)

		case resp := <-resps: // new response received
			if resp == nil {
				break
			}

			// increment
			count++
			data += resp.Size()
			if resp.Err() != nil {
				//fmt.Printf("  %v\n", resp.Err())
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
	var b []byte

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
			}

			// send
			if err == nil {
				hresp, err = client.Do(hreq)
				if err != nil {
					resp.err = err
				}
			}

			// read body
			if err == nil {
				defer hresp.Body.Close()

				resp.resp = hresp
				b, err = ioutil.ReadAll(hresp.Body)
				if err != nil {
					resp.err = err
				}
			}

			if err == nil {
				resp.size = uint64(len(b))
				if req.Bytes > 0 && uint64(len(b)) != req.Bytes {
					resp.err = fmt.Errorf("expected %d bytes; got %d", req.Bytes, resp.size)
				}
			}

		}

		resp.duration = time.Now().Sub(start)
		resps <- resp
	}

	done <- true
}
