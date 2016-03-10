package main

import (
	"time"
)

type ServerStat struct {
	BytesTansferred int64
	Connections     int64
	ConnectionsOpen int64
	Errors          int64
	Requests        int64
	TimeInRequest   time.Duration
}

func (c *ServerStat) Add(n *ServerStat) {
	c.BytesTansferred += n.BytesTansferred
	c.Connections += n.Connections
	c.ConnectionsOpen += n.ConnectionsOpen
	c.Errors += n.Errors
	c.Requests += n.Requests
	c.TimeInRequest += n.TimeInRequest
}
