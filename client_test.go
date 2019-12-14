package main

import (
	"bytes"
	"testing"
)

func TestClient_ReadPump(t *testing.T) {
	client := Client{"test", nil}
	reader := TestReader{[]byte("test buffer")}
	broadcast := make(chan *Message, 1)

	client.ReadPump(&reader, broadcast)

	message := <-broadcast
	if bytes.Compare(reader.read, message.buffer) != 0 {
		t.Errorf("buffer mismatch: %s", string(reader.read))
	}
}

func TestClient_WritePump(t *testing.T) {
	client := Client{"test", make(chan *Message, 1)}
	message := &Message{buffer: []byte("test buffer")}
	w := TestWriter{}

	client.send <- message
	close(client.send)
	client.WritePump(&w, &TestFormatter{})

	if bytes.Compare(w.written, message.buffer) != 0 {
		t.Errorf("buffer mismatch: %s", string(w.written))
	}
}

func TestClient_ReadCallback(t *testing.T) {
	client := Client{"test", nil}
	broadcast := make(chan *Message, 1)
	buffer := []byte("test buffer")
	callback := ReadCallback{&client, broadcast}

	length, _ := callback.Write(buffer)
	if len(buffer) != length {
		t.Errorf("invalid length: %d", length)
	}

	msg := <-broadcast

	if msg.from != &client {
		t.Error("client mismatch")
	}
	if bytes.Compare(msg.buffer, buffer) != 0 {
		t.Error("buffer mismatch")
	}
	if msg.mtype != ClientMsg {
		t.Error("invalid message type")
	}
}
