package storage

import (
	"github.com/syncloud/platform/cli"
	"strconv"
	"strings"
)

type FreeSpaceChecker struct {
	executor cli.Executor
}

func NewFreeSpaceChecker(executor cli.Executor) *FreeSpaceChecker {
	return &FreeSpaceChecker{
		executor: executor,
	}

}

func (f *FreeSpaceChecker) HasFreeSpace(device string) (bool, error) {
	output, err := f.executor.CombinedOutput("parted", device, "unit", "%", "print", "free", "--script", "--machine")
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(output), "\n")
	last := lines[len(lines)-1]
	if !strings.Contains(last, "free") {
		return false, nil
	}
	freeString := strings.Split(last, ":")[3]
	freeString = strings.TrimSuffix(freeString, "%")
	free, err := strconv.ParseFloat(freeString, 64)
	if err != nil {
		return false, err
	}

	return free > ExtendableFreePercent, nil
}
