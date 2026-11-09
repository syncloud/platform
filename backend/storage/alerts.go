package storage

import (
	"github.com/syncloud/platform/storage/model"
	"sort"
)

type AlertsHealth interface {
	Get() (*model.StorageHealth, error)
}

type AlertsDisks interface {
	AvailableDisks() ([]model.Disk, error)
}

type AlertsConfig interface {
	StorageAlertsSeen() []string
	SetStorageAlertsSeen(keys []string)
}

type Alerts struct {
	health AlertsHealth
	disks  AlertsDisks
	config AlertsConfig
}

func NewAlerts(health AlertsHealth, disks AlertsDisks, config AlertsConfig) *Alerts {
	return &Alerts{
		health: health,
		disks:  disks,
		config: config,
	}
}

func (a *Alerts) Get() (*model.StorageAlerts, error) {
	problems, err := a.problems()
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool)
	for _, key := range a.config.StorageAlertsSeen() {
		seen[key] = true
	}
	result := &model.StorageAlerts{Problems: problems}
	for _, problem := range problems {
		if !seen[problem.Key()] {
			result.New = true
		}
	}
	return result, nil
}

func (a *Alerts) Seen() error {
	problems, err := a.problems()
	if err != nil {
		return err
	}
	var keys []string
	for _, problem := range problems {
		keys = append(keys, problem.Key())
	}
	a.config.SetStorageAlertsSeen(keys)
	return nil
}

func (a *Alerts) problems() ([]model.StorageProblem, error) {
	health, err := a.health.Get()
	if err != nil {
		return nil, err
	}
	disks, err := a.disks.AvailableDisks()
	if err != nil {
		return nil, err
	}
	problems := []model.StorageProblem{}
	if health.Scrub != nil && (health.Scrub.Errors != "" || health.Scrub.Uncorrectable > 0) {
		problems = append(problems, model.StorageProblem{Kind: model.ProblemScrubErrors})
	}
	for _, disk := range disks {
		if disk.HasErrors {
			problems = append(problems, model.StorageProblem{Kind: model.ProblemErrorCounters, Device: disk.Device})
		}
	}
	for _, disk := range health.Disks {
		smart := disk.Smart
		if smart.Available && smart.Passed != nil && !*smart.Passed {
			problems = append(problems, model.StorageProblem{Kind: model.ProblemSmartFailed, Device: disk.Device})
		}
		if smart.SelfTest.LastPassed != nil && !*smart.SelfTest.LastPassed {
			problems = append(problems, model.StorageProblem{Kind: model.ProblemSelfTestFailed, Device: disk.Device})
		}
		if disk.Io.Count > 0 {
			problems = append(problems, model.StorageProblem{Kind: model.ProblemIoErrors, Device: disk.Device})
		}
	}
	sort.Slice(problems, func(i, j int) bool { return problems[i].Key() < problems[j].Key() })
	return problems, nil
}
