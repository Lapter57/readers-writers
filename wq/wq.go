package wq

import (
	"container/list"
	"sync"
)

type WaitingQueue struct {
	mu      sync.Mutex
	queue   *list.List
	waiting chan *list.List
}

func NewWaitingQueue() *WaitingQueue {
	queue := &WaitingQueue{
		queue:   list.New(),
		waiting: make(chan *list.List, 1),
	}
	return queue
}

func (c *WaitingQueue) Enqueue(pid int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue.PushBack(pid)
}

func (c *WaitingQueue) Wait(pid int64) {
	for {
		queue := <-c.waiting
		topPid := queue.Front()
		if topPid.Value == pid {
			queue.Remove(topPid)
			break
		} else {
			c.waiting <- queue
		}
	}
}

func (c *WaitingQueue) Dequeue() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.queue.Len() != 0 {
		c.waiting <- c.queue
	}
}

func (c *WaitingQueue) IsQueueEmpty() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.queue.Len() == 0
}

func (c *WaitingQueue) GetQueue() *list.List {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.queue
}
