package main

import (
	"fmt"
	"hash/fnv"
)

type MessageType int

const (
	ClientMsg MessageType = 0
	SystemMsg MessageType = 1
)

type Message struct {
	from   *Client
	buffer []byte
	mtype  MessageType
}

type Formatter interface {
	Format(m *Message, to *Client) []byte
}

type AnsiFormatter struct {
}

func (f AnsiFormatter) Format(m *Message, to *Client) []byte {
	if m.from == to && m.mtype == ClientMsg {
		return []byte{}
	}
	ret := m.buffer

	uname := m.from.username + ": "
	if m.mtype == ClientMsg {
		uname = f.Colorize(uname, UserColor(m.from.username))
	}
	ret = append([]byte(uname), ret...)

	if m.mtype == SystemMsg {
		ret = []byte(f.Colorize(string(ret), SystemColor))
		ret = append(ret, byte('\n'))
	}
	return ret
}

func (f AnsiFormatter) Colorize(s string, c Color) string {
	return fmt.Sprintf("\033[38;05;%dm%s\033[0m", c, s)
}

type Color uint8

var AllColors = ReadableColors()

const SystemColor = Color(245) // grey

func UserColor(username string) Color {
	h := fnv.New32a()
	h.Write([]byte(username))
	return Color(h.Sum32() % uint32(len(AllColors)))
}

func ReadableColors() []Color {
	colors := []Color{}
	var i uint8
	for i = 0; i < 255; i++ {
		if i == 0 || i == 7 || i == 8 || i == 15 || i == 16 || i == 17 || i > 230 {
			continue
		}
		colors = append(colors, Color(i))
	}
	return colors
}
