package main

import (
	"strings"
	"testing"
)

func TestPipe_Run(t *testing.T) {
	pipe := &Pipe{
		broadcast:  make(chan *Message, 1),
		register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
		clients:    make(map[*Client]bool),
	}

	client := &Client{"client1", false, make(chan *Message, 1)}
	message := &Message{client, []byte("test buffer"), SystemMsg}
	go pipe.Run()
	pipe.Register(client)
	joined := <-client.send
	pipe.broadcast <- message
	received := <-client.send
	pipe.Unregister(client)

	if joined.from != client {
		t.Error("invalid join from")
	}
	if joined.mtype != SystemMsg {
		t.Error("invalid join type")
	}
	if !strings.Contains(string(joined.buffer), "joined") {
		t.Error("invalid join buffer")
	}
	if received != message {
		t.Errorf("message mismatch: %s", string(received.buffer))
	}
}
