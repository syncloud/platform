package version

import (
	"os"
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
	content, err := os.ReadFile("/snap/platform/current/META/version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil

}
