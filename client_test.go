package main

import (
	"bytes"
	"testing"
	"time"
)

func TestClient_ReadPump(t *testing.T) {
	client := testClient()
	reader := TestReader{[]byte("test buffer")}
	broadcast := make(chan *Message, 1)

	client.ReadPump(&reader, broadcast)

	message := <-broadcast
	if bytes.Compare(reader.read, message.buffer) != 0 {
		t.Errorf("buffer mismatch: %s", string(reader.read))
	}
}

func TestClient_WritePump(t *testing.T) {
	client := testClient()
	message := &Message{buffer: []byte("test buffer")}
	w := TestWriter{}

	client.send <- message
	close(client.send)
	client.WritePump(&w, make(chan struct{}))

	if bytes.Compare(w.written, message.buffer) != 0 {
		t.Errorf("buffer mismatch: %s", string(w.written))
	}
}

func TestClient_WriteDrip(t *testing.T) {
	client := testClient()
	client.agent = UserAgent{true, 0, 0}
	w := TestWriter{}
	done := make(chan struct{})

	go client.WritePump(&w, done)

	for {
		if w.written != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	var d struct{}
	done <- d

	if bytes.Compare(w.written, []byte(noop)) != 0 {
		t.Errorf("noop not dripped: %s", string(w.written))
	}
}

func TestClient_ReadCallback(t *testing.T) {
	client := testClient()
	broadcast := make(chan *Message, 1)
	buffer := []byte("test buffer")
	callback := ReadCallback{client, broadcast}

	length, _ := callback.Write(buffer)
	if len(buffer) != length {
		t.Errorf("invalid length: %d", length)
	}

	msg := <-broadcast

	if msg.from != client {
		t.Error("client mismatch")
	}
	if bytes.Compare(msg.buffer, buffer) != 0 {
		t.Error("buffer mismatch")
	}
	if msg.mtype != ClientMsg {
		t.Error("invalid message type")
	}
}

func TestClient_IngoreEmpty(t *testing.T) {
	client := testClient()
	broadcast := make(chan *Message, 1)
	buffer := []byte("\n")
	callback := ReadCallback{client, broadcast}
	callback.Write(buffer)

	var broadcasted bool
	select {
	case <-broadcast:
		broadcasted = true
	default:
		broadcasted = false
	}
	if broadcasted {
		t.Error("broadcasted empty message")
	}
}

func testClient() *Client {
	return &Client{"test", UserAgent{false, 0, 0}, make(chan *Message, 1), TestFormatter{}}
}
