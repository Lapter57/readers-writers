package rw

import (
	"rwsync/monitor"
	"sync/atomic"
)

var newWriterId int64 = -1

type Writer struct {
	id         int64
	ob         *Observer
	mon        *monitor.Monitor
	readerCond *monitor.Condition
	writerCond *monitor.Condition
}

func NewWriter(
	ob *Observer,
	mon *monitor.Monitor,
	readerCond *monitor.Condition,
	writerCond *monitor.Condition) *Writer {
	atomic.AddInt64(&newWriterId, 1)
	return &Writer{
		id:         newWriterId,
		ob:         ob,
		mon:        mon,
		readerCond: readerCond,
		writerCond: writerCond,
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
	w.mon.Enter()
	//fmt.Println("START WRITE ", w.id)
	defer w.mon.Leave()

	if w.ob.getReaders() != 0 || w.ob.getWriters() != 0 {
		//fmt.Println("WRITER BLOCKED ", w.id)
		w.mon.Wait(w.writerCond, w.id)
		//fmt.Println("WRITER FREE ", w.id)
	}
	w.ob.AddWriter()
}

func (w *Writer) endWrite() {
	w.mon.Enter()
	//fmt.Println("END WRITE ", w.id)
	defer w.mon.Leave()

	w.ob.RemoveWriter()
	if w.readerCond.IsQueueEmpty() && !w.writerCond.IsQueueEmpty() {
		w.mon.Signal(w.writerCond)
	} else if !w.readerCond.IsQueueEmpty() && w.writerCond.IsQueueEmpty() {
		w.ob.freeReader()
		w.mon.Signal(w.readerCond)
	} else if !w.readerCond.IsQueueEmpty() && !w.writerCond.IsQueueEmpty() {
		w.ob.increaseNumMissWriters()
		reader := w.ob.getWaitingReaders().Front().Value.(*Reader)
		if reader.numMissWriters == 2 {
			w.ob.freeReader()
			w.mon.Signal(w.readerCond)
		} else {
			w.mon.Signal(w.writerCond)
		}
	}
}
