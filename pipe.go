package main

import (
	"sync/atomic"
)

type Pipe struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	maxID      uint64
}

func NewPipe() *Pipe {
	return &Pipe{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (p *Pipe) NextID() uint64 {
	return atomic.AddUint64(&p.maxID, 1)
}

func (p *Pipe) Register(c *Client) {
	p.register <- c
}

func (p *Pipe) Unregister(c *Client) {
	p.unregister <- c
}

func (p *Pipe) Run() {
	for {
		select {
		case client := <-p.register:
			p.add(client)
		case client := <-p.unregister:
			p.remove(client)
		case message := <-p.broadcast:
			p.enqueue(message)
		}
		if len(p.clients) < 1 {
			return
		}
	}
}

func (p *Pipe) add(client *Client) {
	p.clients[client] = true
	p.enqueue(&Message{client, []byte("joined\n"), SystemMsg})
}

func (p *Pipe) remove(client *Client) {
	_, exist := p.clients[client]
	if !exist {
		return
	}
	delete(p.clients, client)
	if client.send != nil {
		close(client.send)
	}
	p.enqueue(&Message{client, []byte("left\n"), SystemMsg})
}

func (p *Pipe) enqueue(message *Message) {
	for client := range p.clients {
		if client.send == nil {
			continue
		}
		select {
		case client.send <- message:
		default:
			p.remove(client) // write queue is full
		}
	}
}
