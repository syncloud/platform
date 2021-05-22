package version

import (
	"io/ioutil"
	"strings"
)

func PlatformVersion() (string, error) {
	content, err := ioutil.ReadFile("/snap/platform/current/META/version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil

}
