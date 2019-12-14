package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Server struct {
	pipe    *Pipe
	baseURL string
}

func (s *Server) GetPipe() *Pipe {
	if s.pipe == nil {
		s.pipe = NewPipe()
		go func() {
			s.pipe.Run()
			s.pipe = nil
		}()
	}
	return s.pipe
}

func (s *Server) Connect(w http.ResponseWriter, r *http.Request) {
	pipe := s.GetPipe()
	client := &Client{getUserName(r, pipe.NextID()), make(chan *Message, 256)}
	defer pipe.Unregister(client)
	pipe.Register(client)

	go client.ReadPump(r.Body, s.pipe.broadcast)
	time.Sleep(10 * time.Millisecond)
	setHeaders(w)
	writer := WriteFlusher{w, getFlusher(w)}
	printWelcome(writer)
	go client.WritePump(writer, AnsiFormatter{})
	<-r.Context().Done()
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func printWelcome(w io.Writer) {
	fmt.Fprintf(w, "Welcome to curlchat\n")
}

func getFlusher(w http.ResponseWriter) http.Flusher {
	flusher, _ := w.(http.Flusher)
	return flusher
}

func getUserName(r *http.Request, id uint64) string {
	username, _, _ := r.BasicAuth()
	if len(username) > 0 {
		return username
	}
	return fmt.Sprintf("user %d", id)
}
