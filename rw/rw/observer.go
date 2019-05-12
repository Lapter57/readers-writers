package rw

import (
	"container/list"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Observer struct {
	mu             sync.Mutex
	rs, ws         int
	waitingReaders *list.List
}

func NewObserver() *Observer {
	return &Observer{
		waitingReaders: list.New(),
	}
}

func (ob *Observer) AddReader() {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.rs++
	ob.printCritSection()
}

func (ob *Observer) RemoveReader() {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.rs--
	ob.printCritSection()
}

func (ob *Observer) AddWriter() {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.ws++
	ob.printCritSection()
}

func (ob *Observer) RemoveWriter() {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.ws--
	ob.printCritSection()
}

func (ob *Observer) printCritSection() {
	fmt.Printf("%s%s\n", strings.Repeat("R", ob.rs),
		strings.Repeat("W", ob.ws))
}

func (ob *Observer) getReaders() int {
	return ob.rs
}

func (ob *Observer) getWriters() int {
	return ob.ws
}

func (ob *Observer) getWaitingReaders() *list.List {
	return ob.waitingReaders
}

func (ob *Observer) freeReader() {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	topReader := ob.waitingReaders.Front()
	ob.waitingReaders.Remove(topReader)
	reader := topReader.Value.(*Reader)
	reader.numMissWriters = 0
}

func (ob *Observer) increaseNumMissWriters() {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	for r := ob.waitingReaders.Front(); r != nil; r = r.Next() {
		reader := r.Value.(*Reader)
		reader.numMissWriters++
	}
}

func sleep() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}
