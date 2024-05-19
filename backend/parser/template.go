package parser

import (
	"errors"
	"fmt"
	"os"
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

	var templates = template.Must(template.ParseGlob(fmt.Sprintf("%s/*", input)))
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
