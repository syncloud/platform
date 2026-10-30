package cron

import "fmt"

type JobStub struct {
	failures int
	runs     int
}

func (j *JobStub) Run() error {
	j.runs++
	if j.runs <= j.failures {
		return fmt.Errorf("network is unreachable")
	}
	return nil
}
