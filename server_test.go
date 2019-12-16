package main

import "testing"

func TestServer_NextID(t *testing.T) {
	server := &Server{}
	id1 := server.NextID()
	if id1 != 1 {
		t.Errorf("inalid nextid: %d", id1)
	}
	id2 := server.NextID()
	if id2 != 2 {
		t.Errorf("inalid nextid: %d", id2)
	}
	id3 := server.NextID()
	if id3 != 3 {
		t.Errorf("inalid nextid: %d", id3)
	}
}
