package version

import (
	"io/ioutil"
)

func PlatformVersion() (string, error) {
	content, err := ioutil.ReadFile("/snap/platform/current/META/version")
	if err != nil {
		return "", err
	}
	return string(content), nil

}
