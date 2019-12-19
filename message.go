package main

import (
	"fmt"
	"hash/fnv"
	"io"
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
	Welcome(w io.Writer, c *Client)
	Message(w io.Writer, m *Message, to *Client)
}

type AnsiFormatter struct {
}

const (
	savecursor    = "\033[s"
	restorecursor = "\033[u"
	insertline    = "\033[L"
	linebeginning = "\033[1G"
	cursordown    = "\033[B"
	colorformat   = "\033[38;05;%dm" // accepts color int
	coloroff      = "\033[0m"
	// Inserts a message above the prompt and doesn't wipe out message in progress
	ansiprefix = savecursor + insertline + linebeginning
	ansisuffix = restorecursor + cursordown
)

func (f AnsiFormatter) Welcome(w io.Writer, c *Client) {
	fmt.Fprintf(w, "Welcome to curlchat\n")
	f.prompt(w, c)
}

func (f AnsiFormatter) Message(w io.Writer, m *Message, to *Client) {
	if m.mtype == ClientMsg && m.from == to {
		f.prompt(w, m.from)
		return
	}

	if m.mtype == ClientMsg {
		fmt.Fprintf(w, ansiprefix+colorformat+"%s: "+coloroff, UserColor(m.from.username), m.from.username)
		w.Write(m.buffer)
		fmt.Fprintf(w, ansisuffix)
		return
	}
	if m.mtype == SystemMsg {
		fmt.Fprintf(w, ansiprefix+colorformat+"%s: ", SystemColor, m.from.username)
		w.Write(m.buffer)
		fmt.Fprintf(w, coloroff+ansisuffix)
		return
	}
}

func (f AnsiFormatter) prompt(w io.Writer, c *Client) {
	fmt.Fprintf(w, colorformat+"%s: "+coloroff, UserColor(c.username), c.username)
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
