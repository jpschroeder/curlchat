package main

import "testing"

func TestPipeCollection_GetPipe(t *testing.T) {
	c := &PipeCollection{}
	p1 := c.GetPipe()

	if p1 == nil {
		t.Error("pipe not initialized")
	}

	p2 := c.GetPipe()

	if p1 != p2 {
		t.Error("not the same pipe")
	}

	client := &Client{}
	p1.Register(client)
	p1.Unregister(client)
	for {
		if c.pipe == nil {
			break
		}
	}
}
