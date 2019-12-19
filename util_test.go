package main

import (
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
	writing []byte
}

func (t *TestWriter) Write(p []byte) (n int, err error) {
	if t.writing == nil {
		t.writing = p
	} else {
		t.writing = append(t.writing, p...)
	}
	return 0, nil
}

func (t *TestWriter) Flush() {
	t.written = t.writing
	t.writing = nil
}

type TestFormatter struct {
}

func (f TestFormatter) Welcome(w io.Writer, c *Client) {
}

func (f TestFormatter) Message(w io.Writer, m *Message, to *Client) {
	w.Write(m.buffer)
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
