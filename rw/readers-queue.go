package rw

import (
	"container/list"
	"sync"
)

type ReadersQueue struct {
	rwMu  sync.RWMutex
	queue *list.List
}

func NewReadersQueue() *ReadersQueue {
	return &ReadersQueue{
		queue: list.New(),
	}
}

func (rq *ReadersQueue) Enqueue(r *Reader) {
	rq.rwMu.Lock()
	defer rq.rwMu.Unlock()
	rq.queue.PushBack(r)
}

func (rq *ReadersQueue) Size() int {
	rq.rwMu.RLock()
	defer rq.rwMu.RUnlock()
	return rq.queue.Len()
}

func (rq *ReadersQueue) Dequeue() {
	rq.rwMu.Lock()
	defer rq.rwMu.Unlock()
	topReader := rq.queue.Front()
	rq.queue.Remove(topReader)
	reader := topReader.Value.(*Reader)
	reader.numMissWriters = 0
}

func (rq *ReadersQueue) Top() *Reader {
	rq.rwMu.RLock()
	defer rq.rwMu.RUnlock()
	return rq.queue.Front().Value.(*Reader)
}

func (rq *ReadersQueue) IncreaseNumMissWriters() {
	rq.rwMu.Lock()
	defer rq.rwMu.Unlock()
	for el := rq.queue.Front(); el != nil; el = el.Next() {
		reader := el.Value.(*Reader)
		reader.numMissWriters++
	}
}
