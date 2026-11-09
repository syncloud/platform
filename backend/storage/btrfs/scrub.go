package btrfs

import (
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"regexp"
	"strconv"
	"strings"
)

var scrubField = regexp.MustCompile(`^\s*([A-Za-z ]+):\s*(.*)$`)

type Scrub struct {
	config   Config
	executor cli.Executor
}

func NewScrub(config Config, executor cli.Executor) *Scrub {
	return &Scrub{
		config:   config,
		executor: executor,
	}
}

func (s *Scrub) Start() error {
	_, err := s.executor.CombinedOutput(BTRFS, "scrub", "start", "-r", s.config.ExternalDiskDir())
	return err
}

func (s *Scrub) Status() (*model.ScrubStatus, error) {
	output, err := s.executor.CombinedOutput(BTRFS, "scrub", "status", s.config.ExternalDiskDir())
	if err != nil {
		return nil, err
	}
	status := ParseScrubStatus(string(output))
	return &status, nil
}

func ParseScrubStatus(output string) model.ScrubStatus {
	status := model.ScrubStatus{State: model.ScrubNever}
	running := false
	for _, line := range strings.Split(output, "\n") {
		match := scrubField.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		key, value := strings.TrimSpace(match[1]), strings.TrimSpace(match[2])
		switch key {
		case "Status":
			status.State = value
		case "Scrub started", "Scrub resumed":
			status.Started = value
		case "Duration":
			status.Duration = value
		case "Time left":
			running = true
		case "Bytes scrubbed":
			running = running || status.State == model.ScrubNever
			if open := strings.Index(value, "("); open >= 0 {
				status.Progress = strings.TrimSuffix(value[open+1:], ")")
			}
		case "Error summary":
			if value != "no errors found" {
				status.Errors = value
			}
		case "Corrected":
			status.Corrected = parseCount(value)
		case "Uncorrectable":
			status.Uncorrectable = parseCount(value)
		}
	}
	if status.State == model.ScrubNever && running {
		status.State = model.ScrubRunning
	}
	return status
}

func parseCount(value string) int64 {
	count, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return count
}
