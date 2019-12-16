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
	out := f.Format(m, c)
	if len(out) > 0 {
		t.Error("shouldn't have echoed client msg")
	}
}

func TestMessage_FormatSystemMessage(t *testing.T) {
	f := AnsiFormatter{}
	c := &Client{}
	m := &Message{c, []byte("test buffer"), SystemMsg}
	out := f.Format(m, c)
	if len(out) < 1 {
		t.Error("should have echoed system msg")
	}
	if !strings.Contains(string(out), strconv.Itoa(int(SystemColor))) {
		t.Errorf("should have been colored with system color: %s", out)
	}
	if out[len(out)-1] != byte('\n') {
		t.Error("should end with a newline")
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
	if !strings.Contains(string(out), "c1") {
		t.Errorf("username not added: %s", string(out))
	}
	if !strings.Contains(string(out), strconv.Itoa(int(UserColor("c1")))) {
		t.Errorf("username not colored: %s", string(out))
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
