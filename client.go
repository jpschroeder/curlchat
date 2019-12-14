package main

import (
	"io"
)

type Client struct {
	username string
	send     chan *Message
}

func (c *Client) ReadPump(reader io.Reader, broadcast chan<- *Message) {
	handler := ReadCallback{c, broadcast}
	// will return on disconnect
	io.Copy(handler, reader)
}

func (c *Client) WritePump(writer io.Writer, formatter Formatter) {
	for {
		message, ok := <-c.send
		if !ok {
			break
		}
		writer.Write(formatter.Format(message, c))
	}
}

type ReadCallback struct {
	from      *Client
	broadcast chan<- *Message
}

func (r ReadCallback) Write(buffer []byte) (int, error) {
	r.broadcast <- &Message{r.from, buffer, ClientMsg}
	return len(buffer), nil
}
