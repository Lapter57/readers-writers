package monitor

import (
	"container/list"
	"sync"
)

type Condition struct {
	mu      sync.Mutex
	queue   *list.List
	waiting chan *list.List
}

func NewCondition() *Condition {
	condition := &Condition{
		queue:   list.New(),
		waiting: make(chan *list.List, 1),
	}
	return condition
}

func (c *Condition) enqueue(pid int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue.PushBack(pid)
}

func (c *Condition) wait(pid int64) {
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

func (c *Condition) dequeue() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.queue.Len() != 0 {
		c.waiting <- c.queue
	}
}

func (c *Condition) IsQueueEmpty() bool {
	return c.queue.Len() == 0
}
