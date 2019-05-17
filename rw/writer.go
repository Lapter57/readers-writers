package rw

import (
	"math/rand"
	"readers-writers/wq"
	"sync"
	"sync/atomic"
	"time"
)

var newWriterId int64 = -1

type Writer struct {
	id               int64
	rwc              *RWCounter
	rwMu             *sync.RWMutex
	readersWaitQueue *wq.WaitingQueue
	writersWaitQueue *wq.WaitingQueue
}

func NewWriter(
	rwc *RWCounter,
	rwMu *sync.RWMutex,
	readersWaitQueue *wq.WaitingQueue,
	writersWaitQueue *wq.WaitingQueue) *Writer {
	atomic.AddInt64(&newWriterId, 1)
	return &Writer{
		id:               newWriterId,
		rwc:              rwc,
		rwMu:             rwMu,
		readersWaitQueue: readersWaitQueue,
		writersWaitQueue: writersWaitQueue,
	}
}

func (w *Writer) Execute() {
	for {
		w.startWrite()
		w.Sleep()
		w.endWrite()
	}
}

func (w *Writer) Sleep() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func (w *Writer) startWrite() {
	w.Sleep()
	w.rwMu.Lock()
	// fmt.Println("START WRITE ", w.id)
	defer w.rwMu.Unlock()

	if w.rwc.GetReaders() != 0 || w.rwc.GetWriters() != 0 {
		// fmt.Println("WRITER BLOCKED ", w.id)
		w.wait()
		// fmt.Println("WRITER FREE ", w.id)
	}
	w.rwc.AddWriter()
}

func (w *Writer) endWrite() {
	w.rwMu.Lock()
	// fmt.Println("END WRITE ", w.id)
	defer w.rwMu.Unlock()

	w.rwc.RemoveWriter()
	if w.readersWaitQueue.IsQueueEmpty() && !w.writersWaitQueue.IsQueueEmpty() {
		w.wakeWriter()
	} else if !w.readersWaitQueue.IsQueueEmpty() && w.writersWaitQueue.IsQueueEmpty() {
		w.wakeReader()
	} else if !w.readersWaitQueue.IsQueueEmpty() && !w.writersWaitQueue.IsQueueEmpty() {
		readersQueue.IncreaseNumMissWriters()
		reader := readersQueue.Top()
		if reader.numMissWriters == 2 {
			w.wakeReader()
		} else {
			w.wakeWriter()
		}
	}
}

func (w *Writer) wait() {
	w.writersWaitQueue.Enqueue(w.id)
	w.rwMu.Unlock()
	w.writersWaitQueue.Wait(w.id)
	w.rwMu.Lock()
}

func (w *Writer) wakeReader() {
	readersQueue.Dequeue()
	if !w.readersWaitQueue.IsQueueEmpty() {
		w.readersWaitQueue.Dequeue()
		w.rwMu.Unlock()
		w.rwMu.Lock()
	}
}

func (w *Writer) wakeWriter() {
	if !w.writersWaitQueue.IsQueueEmpty() {
		w.writersWaitQueue.Dequeue()
		w.rwMu.Unlock()
		w.rwMu.Lock()
	}
}
