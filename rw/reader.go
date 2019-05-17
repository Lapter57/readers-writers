package rw

import (
	"math/rand"
	"readers-writers/wq"
	"sync"
	"sync/atomic"
	"time"
)

var newReaderId int64 = -1
var readersQueue *ReadersQueue = NewReadersQueue()

type Reader struct {
	id               int64
	rwc              *RWCounter
	rwMu             *sync.RWMutex
	readersWaitQueue *wq.WaitingQueue
	writersWaitQueue *wq.WaitingQueue
	numMissWriters   int
}

func NewReader(
	rwc *RWCounter,
	rwMu *sync.RWMutex,
	readersWaitQueue *wq.WaitingQueue,
	writersWaitQueue *wq.WaitingQueue) *Reader {
	atomic.AddInt64(&newReaderId, 1)
	return &Reader{
		id:               newReaderId,
		rwc:              rwc,
		rwMu:             rwMu,
		readersWaitQueue: readersWaitQueue,
		writersWaitQueue: writersWaitQueue,
	}
}

func (r *Reader) Execute() {
	for {
		r.startRead()
		r.Sleep()
		r.endRead()
	}
}

func (r *Reader) Sleep() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func (r *Reader) startRead() {
	r.Sleep()
	r.rwMu.RLock()
	// fmt.Println("START READ ", r.id)
	defer r.rwMu.RUnlock()
	if r.rwc.GetWriters() != 0 || !r.writersWaitQueue.IsQueueEmpty() {
		// fmt.Println("READER BLOCKED ", r.id)
		r.wait()
		// fmt.Println("READER FREE ", r.id)
	}
	r.rwc.AddReader()
	if readersQueue.Size() != 0 {
		reader := readersQueue.Top()
		if reader.numMissWriters == 2 || r.writersWaitQueue.IsQueueEmpty() {
			r.wakeReader()
		}
	}
}

func (r *Reader) endRead() {
	r.rwMu.RLock()
	// fmt.Println("END READ ", r.id)
	defer r.rwMu.RUnlock()
	r.rwc.RemoveReader()
	if r.rwc.GetReaders() == 0 {
		r.wakeWriter()
	}
}

func (r *Reader) wait() {
	readersQueue.Enqueue(r)
	r.readersWaitQueue.Enqueue(r.id)
	r.rwMu.RUnlock()
	r.readersWaitQueue.Wait(r.id)
	r.rwMu.RLock()
}

func (r *Reader) wakeReader() {
	readersQueue.Dequeue()
	if !r.readersWaitQueue.IsQueueEmpty() {
		r.readersWaitQueue.Dequeue()
		r.rwMu.RUnlock()
		r.rwMu.RLock()
	}
}

func (r *Reader) wakeWriter() {
	if !r.writersWaitQueue.IsQueueEmpty() {
		r.writersWaitQueue.Dequeue()
		r.rwMu.RUnlock()
		r.rwMu.RLock()
	}
}
