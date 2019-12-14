package main

import (
	"io"
	"net/http"
)

type WriteFlusher struct {
	writer  io.Writer
	flusher http.Flusher
}

func (wf WriteFlusher) Write(buffer []byte) (int, error) {
	n, e := wf.writer.Write(buffer)
	wf.flusher.Flush()
	return n, e
}
