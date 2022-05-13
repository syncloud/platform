package info

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserConfigMock struct {
	deviceDomain string
	port         int
}

func (u *UserConfigMock) GetPublicPort() *int {
	return &u.port
}

func (u *UserConfigMock) GetDeviceDomain() string {
	return u.deviceDomain
}

func TestUrl_StandardPort(t *testing.T) {
	userConfig := &UserConfigMock{"domain.tld", 443}
	device := New(userConfig)
	url := device.Url("app1")
	assert.Equal(t, "https://app1.domain.tld", url)
}

func TestUrl_NonStandardPort(t *testing.T) {
	userConfig := &UserConfigMock{"domain.tld", 10000}
	device := New(userConfig)
	url := device.Url("app1")
	assert.Equal(t, "https://app1.domain.tld:10000", url)
}
