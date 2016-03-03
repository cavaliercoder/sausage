package main

import (
	"fmt"
	"net/http"
	"time"
)

// ClientRequest represent a single HTTP/S request to be submitted to a proxy
// server. It includes details of the expected response.
type ClientRequest struct {
	LineNo       uint64
	Method       string
	URL          string
	Agent        string
	MimeType     string
	ResponseCode int
	Bytes        uint64
	Duration     time.Duration
	err          error
}

func (c *ClientRequest) Err() error {
	return c.err
}

func (c *ClientRequest) String() string {
	if c.err != nil {
		return fmt.Sprintf("line %d: %s", c.LineNo, c.err.Error())
	}

	return fmt.Sprintf("line %d: %s %s", c.LineNo, c.Method, c.URL)
}

// HTTPRequest returns a http.Request with the ClientRequest encoded into
// X-Sausage headers. These headers are used by sausage server to determine what
// should be returned for the request.
func (c *ClientRequest) HTTPRequest() (*http.Request, error) {
	// create request
	req, err := http.NewRequest(c.Method, c.URL, nil)
	if err != nil {
		return nil, err
	}

	// keep alive on this worker
	req.Header.Set("Connection", "keep-alive")

	// set agent string
	if c.Agent != "" {
		req.Header.Set("User-Agent", c.Agent)
	} else {
		req.Header.Set("User-Agent", "Sausage/1.0")
	}

	// set mimetype
	if c.MimeType != "" {
		req.Header.Set("Accept", c.MimeType)
		req.Header.Set("X-Sausage-Content-Type", c.MimeType)
	}

	// set response ocde
	if c.ResponseCode > 0 {
		req.Header.Set("X-Sausage-Status", fmt.Sprintf("%d", c.ResponseCode))
	}

	// set content length
	if c.Bytes > 0 {
		req.Header.Set("X-Sausage-Content-Length", fmt.Sprintf("%d", c.Bytes))
	}

	// set duration
	if c.Duration > 0 {
		req.Header.Set("X-Sausage-Duration", fmt.Sprintf("%d", c.Duration))
	}

	return req, nil
}
