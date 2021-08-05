package sse

import (
	"github.com/r3labs/sse"
)

// CreateClientConnection Create a connection to the SSE server
func CreateClientConnection(url string) *sse.Client {
	return sse.NewClient(url)
}

// Subscribe subscribes to events and hand them off to a callback function
func Subscribe(client *sse.Client, callback func(msg *sse.Event)) {
	if client == nil {
		panic("invalid sse client connection")
	}

	events := make(chan *sse.Event)
	err := client.SubscribeChan("messages", events)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			msg := <-events
			callback(msg)
		}
	}()
	<-events
}
