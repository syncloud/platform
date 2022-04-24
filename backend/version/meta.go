package version

import (
	"io/ioutil"
	"strings"
)

type PlatformVersion struct{}

type Version interface {
	Get() (string, error)
}

func New() *PlatformVersion {
	return &PlatformVersion{}
}

func (v *PlatformVersion) Get() (string, error) {
	content, err := ioutil.ReadFile("/snap/platform/current/META/version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil

}
