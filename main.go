package main

import (
	"math/rand"
	"readers-writers/pool"
	"readers-writers/rw"
	"readers-writers/wq"
	"sync"
	"time"
)

const (
	numReaders = 3
	numWriters = 3
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var rwc *rw.RWCounter = rw.NewRWCounter()
	var rwMu sync.RWMutex
	var readersQueue *wq.WaitingQueue = wq.NewWaitingQueue()
	var writersQueue *wq.WaitingQueue = wq.NewWaitingQueue()
	var pool *pool.Pool = pool.NewPool(numReaders + numWriters)

	for i := 0; i < numReaders; i++ {
		pool.Exec(rw.NewReader(rwc, &rwMu, readersQueue, writersQueue))
	}

	for i := 0; i < numWriters; i++ {
		pool.Exec(rw.NewWriter(rwc, &rwMu, readersQueue, writersQueue))
	}

	pool.Wait()
}
