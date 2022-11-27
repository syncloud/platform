package backup

import (
	"fmt"
	"regexp"
)

type File struct {
	Path     string `json:"path"`
	File     string `json:"file"`
	App      string `json:"app"`
	FullName string `json:"-"`
}

func Parse(path string, fileName string) (File, error) {
	r, err := regexp.Compile(`(.*?)-\d{4}-\d{4}.*`)
	if err != nil {
		return File{}, err
	}

	matches := r.FindStringSubmatch(fileName)
	if len(matches) < 2 {
		return File{}, fmt.Errorf("backup file name should start with '[app]-YYYY-MMDD-'")
	}
	app := matches[1]
	return File{
		Path:     path,
		File:     fileName,
		App:      app,
		FullName: fmt.Sprintf("%s/%s", path, fileName),
	}, nil
}
