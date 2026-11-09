package model

const (
	ProblemScrubErrors    = "scrub_errors"
	ProblemSelfTestFailed = "self_test_failed"
	ProblemSmartFailed    = "smart_failed"
	ProblemErrorCounters  = "error_counters"
	ProblemIoErrors       = "io_errors"
)

type StorageProblem struct {
	Kind   string `json:"kind"`
	Device string `json:"device,omitempty"`
}

func (p StorageProblem) Key() string {
	return p.Kind + ":" + p.Device
}
