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

func (t *TestReader) Close() error {
	return nil
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

func TestUtil_IsOldCurl(t *testing.T) {
	if isOldCurl("curl/7.58.0") == false {
		t.Error("7.58 is old")
	}
	if isOldCurl("curl/7.67.9") == false {
		t.Error("7.67 is old")
	}
	if isOldCurl("curl/6.68.0") == false {
		t.Error("6.68 is old")
	}
	if isOldCurl("curl/6.70.0") == false {
		t.Error("6.70 is old")
	}
	if isOldCurl("curl/7.67.0-DEV") == false {
		t.Error("7.67 is old")
	}
	if isOldCurl("curl/7.68.0-DEV") == true {
		t.Error("7.68 is new")
	}
	if isOldCurl("curl/7.68.0") == true {
		t.Error("7.68 is new")
	}
	if isOldCurl("curl/7.68.5") == true {
		t.Error("7.68 is new")
	}
	if isOldCurl("curl/8.00.0") == true {
		t.Error("8.00 is new")
	}
	if isOldCurl("blahbloo") == true {
		t.Error("everything else isn't old")
	}
	if isOldCurl("curly/1.00.0") == true {
		t.Error("everything else isn't old")
	}
}
