package sse

import (
	"github.com/r3labs/sse"
	"testing"
)

func TestSubscribe(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	var sseURL = "https://live-test-scores.herokuapp.com/scores"
	client := CreateClientConnection(sseURL)
	Subscribe(client, func(msg *sse.Event) {})

	// Test with an invalid (nil) client
	Subscribe(nil, func(msg *sse.Event) {})
}
