package main

// ClientRequestSource is an interface which returns a channel across which will
// be sent a feed of ClientRequests which may then be sent to a server.
type ClientRequestSource interface {
	Get() (<-chan *ClientRequest, error)
}
