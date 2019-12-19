package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"
)

const noop = "\033[0m" // use reset color as noop

type Client struct {
	username string
	oldcurl  bool
	send     chan *Message
}

func (c *Client) ReadPump(reader io.ReadCloser, broadcast chan<- *Message) {
	log.Printf("%s: reader started", c.username)
	handler := ReadCallback{c, broadcast}
	// will return on disconnect
	io.Copy(handler, reader)
	reader.Close()
	log.Printf("%s: reader done", c.username)
}

func (c *Client) WritePump(writer WriteFlusher, done <-chan struct{}, formatter Formatter) {
	log.Printf("%s: writer started", c.username)

	formatter.Welcome(writer, c)
	writer.Flush()

	// drip a stream of noop characters to fix a bug in old versions of curl
	ticker := time.NewTicker(500 * time.Millisecond)
	if !c.oldcurl {
		ticker.Stop()
	} else {
		defer ticker.Stop()
	}

writeloop:
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				break writeloop
			}
			formatter.Message(writer, message, c)
			writer.Flush()
		case <-ticker.C:
			fmt.Fprintf(writer, noop) // reset color as noop
			writer.Flush()
		case <-done:
			break writeloop
		}
	}
	log.Printf("%s: writer done", c.username)
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
