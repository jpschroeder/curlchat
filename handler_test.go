package main

import "testing"

func TestHandler_NextID(t *testing.T) {
	handler := &Handler{}
	id1 := handler.NextID()
	if id1 != 1 {
		t.Errorf("inalid nextid: %d", id1)
	}
	id2 := handler.NextID()
	if id2 != 2 {
		t.Errorf("inalid nextid: %d", id2)
	}
	id3 := handler.NextID()
	if id3 != 3 {
		t.Errorf("inalid nextid: %d", id3)
	}
}
