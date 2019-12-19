package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Server struct {
	pipes *PipeCollection
	maxID uint64
}

func (s *Server) Connect(w http.ResponseWriter, r *http.Request) {
	pipe := s.pipes.GetPipe()
	client := &Client{
		username: getUserName(r, s.NextID()),
		oldcurl:  isOldCurl(r.UserAgent()),
		send:     make(chan *Message, 256)}
	defer pipe.Unregister(client)
	pipe.Register(client)

	go client.ReadPump(r.Body, pipe.broadcast)
	time.Sleep(10 * time.Millisecond)
	setHeaders(w)
	client.WritePump(WriteFlush{w, getFlusher(w)}, r.Context().Done(), AnsiFormatter{})
}

func (s *Server) NextID() uint64 {
	return atomic.AddUint64(&s.maxID, 1)
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
