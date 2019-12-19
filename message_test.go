package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestMessage_FormatNoEchoClient(t *testing.T) {
	f := AnsiFormatter{}
	c := &Client{}
	m := &Message{c, []byte("test buffer"), ClientMsg}
	w := &TestWriter{}
	f.Message(w, m, c)
	w.Flush()
	if strings.Contains(string(w.written), "test buffer") {
		t.Errorf("shouldn't have echoed client msg: %s", w.written)
	}
}

func TestMessage_FormatSystemMessage(t *testing.T) {
	f := AnsiFormatter{}
	c := &Client{}
	m := &Message{c, []byte("test buffer"), SystemMsg}
	w := &TestWriter{}
	f.Message(w, m, c)
	w.Flush()
	if len(w.written) < 1 {
		t.Errorf("should have echoed system msg: %s", w.written)
	}
	if !strings.Contains(string(w.written), strconv.Itoa(int(SystemColor))) {
		t.Errorf("should have been colored with system color: %s", w.written)
	}
}

func TestMessage_FormatUsername(t *testing.T) {
	f := AnsiFormatter{}
	c1 := &Client{username: "c1"}
	c2 := &Client{username: "c2"}
	m := &Message{c1, []byte("test buffer"), ClientMsg}
	w := &TestWriter{}
	f.Message(w, m, c2)
	w.Flush()
	if len(w.written) < 1 {
		t.Error("should have sent msg")
	}
	if !strings.Contains(string(w.written), "c1") {
		t.Errorf("username not added: %s", string(w.written))
	}
	if !strings.Contains(string(w.written), strconv.Itoa(int(UserColor("c1")))) {
		t.Errorf("username not colored: %s", string(w.written))
	}
}

func TestMessage_UserColor(t *testing.T) {
	c1 := UserColor("user 1")
	c2 := UserColor("user 1")
	c3 := UserColor("user 2")

	if c1 != c2 {
		t.Error("usernames should always get the same color")
	}
	if c1 == c3 {
		t.Error("different usernames should get different colors")
	}
}

func TestMessage_AllColors(t *testing.T) {
	/*
		f := AnsiFormatter{}
		for _, color := range AllColors {
			fmt.Printf("%d\t%s\n", color, f.Colorize("ABCDEFG", color))
		}
	*/
}
