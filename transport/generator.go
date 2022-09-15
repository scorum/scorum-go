package transport

import "sync"

type SequenceGenerator struct {
	requestID uint64
	mtx       sync.Mutex
}

func NewSequenceGenerator(initial uint64) *SequenceGenerator {
	return &SequenceGenerator{
		requestID: initial,
		mtx:       sync.Mutex{},
	}
}

func (r *SequenceGenerator) Generate() uint64 {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.requestID++

	return r.requestID
}
