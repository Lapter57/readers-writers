package rw

import (
	"fmt"
	"strings"
	"sync"
)

type RWCounter struct {
	mu     sync.Mutex
	rs, ws int
}

func NewRWCounter() *RWCounter {
	return &RWCounter{}
}

func (rwc *RWCounter) AddReader() {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	rwc.rs++
	rwc.printNumRW()
}

func (rwc *RWCounter) RemoveReader() {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	rwc.rs--
	rwc.printNumRW()
}

func (rwc *RWCounter) AddWriter() {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	rwc.ws++
	rwc.printNumRW()
}

func (rwc *RWCounter) RemoveWriter() {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	rwc.ws--
	rwc.printNumRW()
}

func (rwc *RWCounter) printNumRW() {
	fmt.Printf("%s%s\n", strings.Repeat("R", rwc.rs),
		strings.Repeat("W", rwc.ws))
}

func (rwc *RWCounter) GetReaders() int {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	return rwc.rs
}

func (rwc *RWCounter) GetWriters() int {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	return rwc.ws
}
