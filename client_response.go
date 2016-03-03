package main

import (
	"fmt"
	"net/http"
	"time"
)

type ClientResponse struct {
	req      *ClientRequest
	resp     *http.Response
	size     uint64
	duration time.Duration
	err      error
}

func (c *ClientResponse) String() string {
	if c.err != nil {
		return fmt.Sprintf("Error: %v (%d bytes in %v)", c.err, c.size, c.duration)
	} else {
		return fmt.Sprintf("%s (%d bytes in %v)", c.resp.Status, c.size, c.duration)
	}
}

func (c *ClientResponse) Request() *ClientRequest {
	return c.req
}

func (c *ClientResponse) HTTPResponse() *http.Response {
	return c.resp
}

func (c *ClientResponse) Size() uint64 {
	return c.size
}

func (c *ClientResponse) Duration() time.Duration {
	return c.duration
}

func (c *ClientResponse) Err() error {
	return c.err
}
