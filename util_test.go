package main

import (
	"bytes"
	"io"
	"testing"
)

type TestReader struct {
	read []byte
}

func (t *TestReader) Read(p []byte) (n int, err error) {
	copy(p, t.read)
	return len(t.read), io.EOF
}

type TestWriter struct {
	written []byte
}

func (t *TestWriter) Write(p []byte) (n int, err error) {
	t.written = p
	return 0, nil
}

type TestFlusher struct {
	called bool
}

func (t *TestFlusher) Flush() {
	t.called = true
}

type TestFormatter struct {
}

func (f TestFormatter) Format(m *Message, to *Client) []byte {
	return m.buffer
}

func TestUtil_WriteFlusher(t *testing.T) {
	w := TestWriter{}
	f := TestFlusher{}
	wf := WriteFlusher{&w, &f}
	p := []byte("test message")
	wf.Write(p)

	if bytes.Compare(w.written, p) != 0 {
		t.Error("buffer mismatch")
	}

	if !f.called {
		t.Error("flush not called")
	}
}
