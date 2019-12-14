package main

import "testing"

import "bytes"

func TestPipe_Run(t *testing.T) {
	pipe := &Pipe{
		broadcast:  make(chan *Message, 1),
		register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
		clients:    make(map[*Client]bool),
	}

	client := &Client{"client1", make(chan *Message, 1)}
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
	if bytes.Compare(joined.buffer, []byte("joined\n")) != 0 {
		t.Error("invalid join buffer")
	}
	if received != message {
		t.Errorf("message mismatch: %s", string(received.buffer))
	}
}

func TestPipe_NextID(t *testing.T) {
	pipe := &Pipe{}
	id1 := pipe.NextID()
	if id1 != 1 {
		t.Errorf("inalid nextid: %d", id1)
	}
	id2 := pipe.NextID()
	if id2 != 2 {
		t.Errorf("inalid nextid: %d", id2)
	}
	id3 := pipe.NextID()
	if id3 != 3 {
		t.Errorf("inalid nextid: %d", id3)
	}
}
