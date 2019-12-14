package main

import (
	"strings"
	"testing"
)

func TestMessage_FormatNoEchoClient(t *testing.T) {
	f := AnsiFormatter{}
	c := &Client{}
	m := &Message{c, []byte("test buffer"), ClientMsg}
	out := f.Format(m, c)
	if len(out) > 0 {
		t.Error("shouldn't have echoed client msg")
	}
}

func TestMessage_FormatEchoSystem(t *testing.T) {
	f := AnsiFormatter{}
	c := &Client{}
	m := &Message{c, []byte("test buffer"), SystemMsg}
	out := f.Format(m, c)
	if len(out) < 1 {
		t.Error("should have echoed system msg")
	}
}

func TestMessage_FormatUsername(t *testing.T) {
	f := AnsiFormatter{}
	c1 := &Client{username: "c1"}
	c2 := &Client{username: "c2"}
	m := &Message{c1, []byte("test buffer"), ClientMsg}
	out := f.Format(m, c2)
	if len(out) < 1 {
		t.Error("should have sent msg")
	}
	if !strings.HasPrefix(string(out), "c1") {
		t.Errorf("username not added: %s", string(out))
	}
}
