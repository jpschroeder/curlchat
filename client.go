package main

import (
	"bytes"
	"io"
	"time"
)

const noop = "\033[0m" // use reset color as noop

type Client struct {
	username string
	oldcurl  bool
	send     chan *Message
}

func (c *Client) ReadPump(reader io.Reader, broadcast chan<- *Message) {
	handler := ReadCallback{c, broadcast}
	// will return on disconnect
	io.Copy(handler, reader)
}

func (c *Client) WritePump(writer io.Writer, formatter Formatter) {

	// drip a stream of noop characters to fix a bug in old versions of curl
	ticker := time.NewTicker(500 * time.Millisecond)
	if !c.oldcurl {
		ticker.Stop()
	} else {
		defer ticker.Stop()
	}

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			writer.Write(formatter.Format(message, c))
		case <-ticker.C:
			writer.Write([]byte(noop)) // reset color as noop
		}
	}
}

type ReadCallback struct {
	from      *Client
	broadcast chan<- *Message
}

func (r ReadCallback) Write(buffer []byte) (int, error) {
	if bytes.Compare([]byte("\n"), buffer) != 0 {
		r.broadcast <- &Message{r.from, buffer, ClientMsg}
	}
	return len(buffer), nil
}
