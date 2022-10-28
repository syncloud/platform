package job

import (
	"fmt"
	"sync"
)

type SingleJobMaster struct {
	mutex  *sync.Mutex
	status int
	job    func()
	name   string
}

func NewMaster() *SingleJobMaster {
	return &SingleJobMaster{
		mutex:  &sync.Mutex{},
		status: Idle,
	}
}

func (m *SingleJobMaster) Status() Status {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return NewStatus(m.name, m.status)
}

func (m *SingleJobMaster) Offer(name string, job func()) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.status == Idle {
		m.status = Waiting
		m.job = job
		m.name = name
		return nil
	} else {
		return fmt.Errorf("busy")
	}
}

func (m *SingleJobMaster) Take() (func(), error) {

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.status == Waiting {
		m.status = Busy
		return m.job, nil
	} else {
		return nil, fmt.Errorf("no tasks")
	}
}

func (m *SingleJobMaster) Complete() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.status == Busy {
		m.status = Idle
		m.job = nil
		m.name = ""
		return nil
	} else {
		return fmt.Errorf("nothing to complete")
	}
}
