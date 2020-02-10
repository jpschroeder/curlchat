package main

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

const noop = "\033[0m" // use reset color as noop

type Client struct {
	username  string
	agent     UserAgent
	send      chan *Message
	formatter Formatter
}

func (c *Client) ReadPump(reader io.ReadCloser, broadcast chan<- *Message) {
	handler := ReadCallback{c, broadcast}
	// will return on disconnect
	io.Copy(handler, reader)
	reader.Close()
}

func (c *Client) WritePump(writer WriteFlusher, done <-chan struct{}) {
	c.formatter.Welcome(writer, c)
	writer.Flush()

	// drip a stream of noop characters to fix a bug in old versions of curl
	ticker := time.NewTicker(500 * time.Millisecond)
	if !c.agent.isOldCurl() {
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
			c.formatter.Message(writer, message, c)
			writer.Flush()
		case <-ticker.C:
			fmt.Fprintf(writer, noop) // reset color as noop
			writer.Flush()
		case <-done:
			return
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
