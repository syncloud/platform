package job

import (
	"sync"
)

type Job struct {
	
}

type Master struct {
	mutex  *sync.Mutex
	status string
}

func NewMaster() *Master {
	return &Master{
		mutex:  &sync.Mutex{},
		status: "empty",
	}
}

func (master *Master) Status() string {
	master.mutex.Lock()
	defer master.mutex.Unlock()
	return master.status
}
