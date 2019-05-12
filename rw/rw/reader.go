package rw

import (
	"rwsync/monitor"
	"sync/atomic"
)

var newReaderId int64 = -1

type Reader struct {
	id             int64
	ob             *Observer
	mon            *monitor.Monitor
	readerCond     *monitor.Condition
	writerCond     *monitor.Condition
	numMissWriters int
}

func NewReader(
	ob *Observer,
	mon *monitor.Monitor,
	readerCond *monitor.Condition,
	writerCond *monitor.Condition) *Reader {
	atomic.AddInt64(&newReaderId, 1)
	return &Reader{
		id:         newReaderId,
		ob:         ob,
		mon:        mon,
		readerCond: readerCond,
		writerCond: writerCond,
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
	r.mon.Enter()
	//fmt.Println("START READ ", r.id)
	defer r.mon.Leave()
	if r.ob.getWriters() != 0 || !r.writerCond.IsQueueEmpty() {
		r.ob.getWaitingReaders().PushBack(r)
		//fmt.Println("READER BLOCKED ", r.id)
		r.mon.Wait(r.readerCond, r.id)
		//fmt.Println("READER FREE ", r.id)
	}
	r.ob.AddReader()
	if r.ob.getWaitingReaders().Len() != 0 {
		reader := r.ob.getWaitingReaders().Front().Value.(*Reader)
		if reader.numMissWriters == 2 || r.writerCond.IsQueueEmpty() {
			r.ob.freeReader()
			r.mon.Signal(r.readerCond)
		}
	}
}

func (r *Reader) endRead() {
	r.mon.Enter()
	//fmt.Println("END READ ", r.id)
	defer r.mon.Leave()
	r.ob.RemoveReader()
	if r.ob.getReaders() == 0 {
		r.mon.Signal(r.writerCond)
	}
}
