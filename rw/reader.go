package rw

import (
	"container/list"
	"readers-writers/wq"
	"sync"
	"sync/atomic"
)

var newReaderId int64 = -1
var waitingReaders *list.List = list.New()
var mu sync.Mutex

type Reader struct {
	id             int64
	rwc            *RWCounter
	rwMu           *sync.RWMutex
	readersQueue   *wq.WaitingQueue
	writersQueue   *wq.WaitingQueue
	numMissWriters int
}

func NewReader(
	rwc *RWCounter,
	rwMu *sync.RWMutex,
	readersQueue *wq.WaitingQueue,
	writersQueue *wq.WaitingQueue) *Reader {
	atomic.AddInt64(&newReaderId, 1)
	return &Reader{
		id:           newReaderId,
		rwc:          rwc,
		rwMu:         rwMu,
		readersQueue: readersQueue,
		writersQueue: writersQueue,
	}
}

func (r *Reader) Execute() {
	for {
		r.startRead()
		sleep()
		r.endRead()
	}
}

func (r *Reader) startRead() {
	sleep()
	r.rwMu.RLock()
	// fmt.Println("START READ ", r.id)
	defer r.rwMu.RUnlock()
	if r.rwc.GetWriters() != 0 || !r.writersQueue.IsQueueEmpty() {
		// fmt.Println("READER BLOCKED ", r.id)
		r.wait()
		// fmt.Println("READER FREE ", r.id)
	}
	r.rwc.AddReader()
	if waitingReaders.Len() != 0 {
		reader := waitingReaders.Front().Value.(*Reader)
		if reader.numMissWriters == 2 || r.writersQueue.IsQueueEmpty() {
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
	waitingReaders.PushBack(r)
	r.readersQueue.Enqueue(r.id)
	r.rwMu.RUnlock()
	r.readersQueue.Wait(r.id)
	r.rwMu.RLock()
}

func (r *Reader) wakeReader() {
	freeReader()
	if !r.readersQueue.IsQueueEmpty() {
		r.readersQueue.Dequeue()
		r.rwMu.RUnlock()
		r.rwMu.RLock()
	}
}

func (r *Reader) wakeWriter() {
	if !r.writersQueue.IsQueueEmpty() {
		r.writersQueue.Dequeue()
		r.rwMu.RUnlock()
		r.rwMu.RLock()
	}
}

func freeReader() {
	mu.Lock()
	defer mu.Unlock()
	topReader := waitingReaders.Front()
	waitingReaders.Remove(topReader)
	reader := topReader.Value.(*Reader)
	reader.numMissWriters = 0
}

func increaseNumMissWriters() {
	mu.Lock()
	defer mu.Unlock()
	for el := waitingReaders.Front(); el != nil; el = el.Next() {
		reader := el.Value.(*Reader)
		reader.numMissWriters++
	}
}
