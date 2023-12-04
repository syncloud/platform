package model

import "regexp"

type Change struct {
	Id      string `json:"id"`
	Summary string `json:"summary"`
	Tasks   []Task `json:"tasks"`
}

func (c Change) InstallerProgress() InstallerProgress {
	app := ParseApp(c.Summary)
	for _, task := range c.Tasks {
		if task.Status == "Doing" {
			if task.Kind == "download-snap" {
				return InstallerProgress{
					App:           app,
					Summary:       "Downloading",
					Indeterminate: false,
					Percentage:    CalcPercentage(task.Progress.Done, task.Progress.Total),
				}
			}
		}
	}
	return InstallerProgress{
		App:           app,
		Summary:       ParseAction(c.Summary),
		Indeterminate: true,
	}
}

func CalcPercentage(done, total int64) int64 {
	if total <= 1 {
		return 0
	}
	return done * 100 / total
}

func ParseApp(summary string) string {
	r := regexp.MustCompile(`^.*? "(.*?)" .*`)
	match := r.FindStringSubmatch(summary)
	if match != nil {
		return match[1]
	}
	return "unknown"
}

func ParseAction(summary string) string {
	r := regexp.MustCompile(`^(.*?) .*`)
	match := r.FindStringSubmatch(summary)
	if match != nil {
		switch match[1] {
		case "Refresh":
			return "Upgrading"
		case "Install":
			return "Installing"
		case "Remove":
			return "Removing"
		}
	}
	return "Unknown"
}

type Task struct {
	Kind     string   `json:"kind"`
	Status   string   `json:"status"`
	Summary  string   `json:"summary"`
	Progress Progress `json:"progress"`
}

type Progress struct {
	Done  int64 `json:"done"`
	Total int64 `json:"total"`
}
