package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func Generate(input, output string, data interface{}) error {
	_, err := os.Stat(output)
	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(output, 0755)
		if err != nil {
			return err
		}
	}

	entries, err := filepath.Glob(fmt.Sprintf("%s/*", input))
	if err != nil {
		return err
	}
	var files []string
	for _, entry := range entries {
		info, err := os.Stat(entry)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, entry)
		}
	}
	if len(files) == 0 {
		return fmt.Errorf("no template files found in %s", input)
	}

	var templates = template.Must(template.ParseFiles(files...))
	for _, t := range templates.Templates() {
		err := write(output, t, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func write(output string, t *template.Template, data interface{}) error {
	f, err := os.Create(fmt.Sprintf("%s/%s", output, t.Name()))
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	return nil
}
