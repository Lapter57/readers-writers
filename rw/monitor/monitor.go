package monitor

type Monitor struct {
	lock chan struct{}
}

func NewMonitor() *Monitor {
	monitor := &Monitor{
		lock: make(chan struct{}, 1),
	}
	monitor.lock <- struct{}{}
	return monitor
}

func (m *Monitor) Enter() {
	<-m.lock
}

func (m *Monitor) Leave() {
	m.lock <- struct{}{}
}

func (m *Monitor) Wait(c *Condition, pid int64) {
	c.enqueue(pid)
	m.lock <- struct{}{}
	c.wait(pid)
	if len(m.lock) == 1 {
		<-m.lock
	}
}

func (m *Monitor) Signal(c *Condition) {
	if !c.IsQueueEmpty() {
		c.dequeue()
		<-m.lock
	}
}
