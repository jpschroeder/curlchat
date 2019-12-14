package main

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
	if len(m.from.username) > 0 {
		ret = append([]byte(m.from.username+": "), ret...)
	}
	return ret
}
