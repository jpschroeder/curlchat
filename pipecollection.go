package main

import "sync"

import "sync/atomic"

type PipeCollection struct {
	pipe        *Pipe
	initialized uint32
	mux         sync.Mutex
}

func (p *PipeCollection) GetPipe() *Pipe {
	if atomic.LoadUint32(&p.initialized) == 1 {
		return p.pipe
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	if p.initialized == 0 {
		p.pipe = NewPipe()
		go func() {
			p.pipe.Run()
			p.pipe = nil
			atomic.StoreUint32(&p.initialized, 0)
		}()
		atomic.StoreUint32(&p.initialized, 1)
	}
	return p.pipe
}
