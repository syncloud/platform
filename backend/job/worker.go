package job

import (
	"log"
	"time"
)

type Worker struct {
	master JobMaster
}

func NewWorker(master JobMaster) *Worker {
	return &Worker{master}
}

func (worker *Worker) Start() {
	for {
		if !worker.Do() {
			time.Sleep(1 * time.Second)
		}
	}
}

func (worker *Worker) Do() bool {
	job, err := worker.master.Take()
	if err != nil {
		return false
	}
	job()
	err = worker.master.Complete()
	if err != nil {
		log.Println("error: ", err)
	}
	return true
}
