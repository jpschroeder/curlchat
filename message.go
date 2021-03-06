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

func getFormatter(agent UserAgent, baseURL string) Formatter {
	if agent.isCurl {
		return AnsiFormatter{baseURL}
	}
	return TextFormatter{baseURL}
}

type AnsiFormatter struct {
	baseURL string
}

const (
	linebeginning = "\033[1G"
	colorformat   = "\033[38;05;%dm" // accepts color int
	coloroff      = "\033[0m"
)

func (f AnsiFormatter) Welcome(w io.Writer, c *Client) {
	fmt.Fprintf(w, "Welcome to curlchat\n")
	fmt.Fprintf(w, "curl -T. -N %s -u username:\n", f.baseURL)
	f.prompt(w, c)
}

func (f AnsiFormatter) Message(w io.Writer, m *Message, to *Client) {
	if m.mtype == ClientMsg && m.from == to {
		f.prompt(w, m.from)
		return
	}

	if m.mtype == ClientMsg {
		fmt.Fprintf(w, linebeginning+colorformat+"%s: "+coloroff, UserColor(m.from.username), m.from.username)
		w.Write(m.buffer)
	}
	if m.mtype == SystemMsg {
		fmt.Fprintf(w, linebeginning+colorformat+"%s: ", SystemColor, m.from.username)
		w.Write(m.buffer)
		fmt.Fprintf(w, "%s\n", coloroff)
	}
	f.prompt(w, to)
}

func (f AnsiFormatter) prompt(w io.Writer, c *Client) {
	fmt.Fprintf(w, colorformat+"%s: "+coloroff, UserColor(c.username), c.username)
}

type TextFormatter struct {
	baseURL string
}

func (f TextFormatter) Welcome(w io.Writer, c *Client) {
	fmt.Fprintf(w, "Welcome to curlchat\n")
	fmt.Fprintf(w, "curl -T. -N %s -u username:\n", f.baseURL)
}

func (f TextFormatter) Message(w io.Writer, m *Message, to *Client) {
	if m.mtype == ClientMsg && m.from == to {
		return
	}

	fmt.Fprintf(w, "%s: ", m.from.username)
	w.Write(m.buffer)
	if m.mtype == SystemMsg {
		w.Write([]byte("\n"))
	}
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
