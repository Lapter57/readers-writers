package main

import (
	"math/rand"
	"rwsync/monitor"
	"rwsync/pool"
	"rwsync/rw"
	"time"
)

const N_READERS int = 3
const N_WRITERS int = 3

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var ob *rw.Observer = rw.NewObserver()
	var mon *monitor.Monitor = monitor.NewMonitor()
	var readerCond *monitor.Condition = monitor.NewCondition()
	var writerCond *monitor.Condition = monitor.NewCondition()

	var pool *pool.Pool = pool.NewPool(N_READERS + N_WRITERS)

	for i := 0; i < N_READERS; i++ {
		pool.Exec(rw.NewReader(ob, mon, readerCond, writerCond))
	}

	for i := 0; i < N_WRITERS; i++ {
		pool.Exec(rw.NewWriter(ob, mon, readerCond, writerCond))
	}

	pool.Wait()
}
