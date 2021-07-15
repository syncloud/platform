package connection

import (
	"fmt"
	"net/http"
)

type Checker interface {
	Check() error
}

type Internet struct {
}

func (i *Internet) Check() error {
	url := "http://apps.syncloud.org/releases/stable/index"
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("internet check url %s is not reachable, error: %s", url, err)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("internet check, response status_code: %d", response.StatusCode)
	}
	return nil
}
