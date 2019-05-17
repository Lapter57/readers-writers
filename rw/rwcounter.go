package rw

import (
	"fmt"
	"strings"
	"sync"
)

type RWCounter struct {
	rwMu   sync.RWMutex
	rs, ws int
}

func NewRWCounter() *RWCounter {
	return &RWCounter{}
}

func (rwc *RWCounter) AddReader() {
	rwc.rwMu.Lock()
	defer rwc.rwMu.Unlock()
	rwc.rs++
	rwc.printNumRW()
}

func (rwc *RWCounter) RemoveReader() {
	rwc.rwMu.Lock()
	defer rwc.rwMu.Unlock()
	rwc.rs--
	rwc.printNumRW()
}

func (rwc *RWCounter) AddWriter() {
	rwc.rwMu.Lock()
	defer rwc.rwMu.Unlock()
	rwc.ws++
	rwc.printNumRW()
}

func (rwc *RWCounter) RemoveWriter() {
	rwc.rwMu.Lock()
	defer rwc.rwMu.Unlock()
	rwc.ws--
	rwc.printNumRW()
}

func (rwc *RWCounter) printNumRW() {
	fmt.Printf("%s%s\n", strings.Repeat("R", rwc.rs),
		strings.Repeat("W", rwc.ws))
}

func (rwc *RWCounter) GetReaders() int {
	rwc.rwMu.RLock()
	defer rwc.rwMu.RUnlock()
	return rwc.rs
}

func (rwc *RWCounter) GetWriters() int {
	rwc.rwMu.RLock()
	defer rwc.rwMu.RUnlock()
	return rwc.ws
}
