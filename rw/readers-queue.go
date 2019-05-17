package rw

import (
	"container/list"
	"sync"
)

type ReadersQueue struct {
	mu    sync.Mutex
	queue *list.List
}

func NewReadersQueue() *ReadersQueue {
	return &ReadersQueue{
		queue: list.New(),
	}
}

func (rq *ReadersQueue) Enqueue(r *Reader) {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	rq.queue.PushBack(r)
}

func (rq *ReadersQueue) Size() int {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	return rq.queue.Len()
}

func (rq *ReadersQueue) Dequeue() {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	topReader := rq.queue.Front()
	rq.queue.Remove(topReader)
	reader := topReader.Value.(*Reader)
	reader.numMissWriters = 0
}

func (rq *ReadersQueue) Top() *Reader {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	return rq.queue.Front().Value.(*Reader)
}

func (rq *ReadersQueue) IncreaseNumMissWriters() {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	for el := rq.queue.Front(); el != nil; el = el.Next() {
		reader := el.Value.(*Reader)
		reader.numMissWriters++
	}
}
