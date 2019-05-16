package rw

import (
	"readers-writers/wq"
	"sync"
	"sync/atomic"
)

var newWriterId int64 = -1

type Writer struct {
	id           int64
	rwc          *RWCounter
	rwMu         *sync.RWMutex
	readersQueue *wq.WaitingQueue
	writersQueue *wq.WaitingQueue
}

func NewWriter(
	rwc *RWCounter,
	rwMu *sync.RWMutex,
	readersQueue *wq.WaitingQueue,
	writersQueue *wq.WaitingQueue) *Writer {
	atomic.AddInt64(&newWriterId, 1)
	return &Writer{
		id:           newWriterId,
		rwc:          rwc,
		rwMu:         rwMu,
		readersQueue: readersQueue,
		writersQueue: writersQueue,
	}
}

func (w *Writer) Execute() {
	for {
		w.startWrite()
		sleep()
		w.endWrite()
	}
}

func (w *Writer) startWrite() {
	sleep()
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
	if w.readersQueue.IsQueueEmpty() && !w.writersQueue.IsQueueEmpty() {
		w.wakeWriter()
	} else if !w.readersQueue.IsQueueEmpty() && w.writersQueue.IsQueueEmpty() {
		w.wakeReader()
	} else if !w.readersQueue.IsQueueEmpty() && !w.writersQueue.IsQueueEmpty() {
		increaseNumMissWriters()
		reader := waitingReaders.Front().Value.(*Reader)
		if reader.numMissWriters == 2 {
			w.wakeReader()
		} else {
			w.wakeWriter()
		}
	}
}

func (w *Writer) wait() {
	w.writersQueue.Enqueue(w.id)
	w.rwMu.Unlock()
	w.writersQueue.Wait(w.id)
	w.rwMu.Lock()
}

func (w *Writer) wakeReader() {
	freeReader()
	if !w.readersQueue.IsQueueEmpty() {
		w.readersQueue.Dequeue()
		w.rwMu.Unlock()
		w.rwMu.Lock()
	}
}

func (w *Writer) wakeWriter() {
	if !w.writersQueue.IsQueueEmpty() {
		w.writersQueue.Dequeue()
		w.rwMu.Unlock()
		w.rwMu.Lock()
	}
}
