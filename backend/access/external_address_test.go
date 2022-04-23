package access

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"io/ioutil"
	"net/http"
	"testing"
)

type UserConfigStub struct {
}

func (u UserConfigStub) GetRedirectApiUrl() string {
	return "url"
}

func (u UserConfigStub) GetDomainUpdateToken() *string {
	token := "token"
	return &token
}

func (u UserConfigStub) IsRedirectEnabled() bool {
	return true
}

func (u UserConfigStub) SetIpv4Enabled(_ bool) {}

func (u UserConfigStub) SetIpv4Public(_ bool) {}

func (u UserConfigStub) SetIpv6Enabled(_ bool) {}

func (u UserConfigStub) SetPublicIp(_ *string) {}

func (u UserConfigStub) SetPublicPort(_ *int) {}

func (u UserConfigStub) GetPublicIp() *string {
	ip := "1.1.1.1"
	return &ip
}

func (u UserConfigStub) GetPublicPort() *int {
	port := 1
	return &port
}

func (u UserConfigStub) IsIpv6Enabled() bool {
	return true
}

func (u UserConfigStub) IsIpv4Public() bool {
	return true
}

func (u UserConfigStub) IsIpv4Enabled() bool {
	return true
}

type RedirectStub struct {
}

func (r RedirectStub) Update(_ *string, _ *string, _ *int, _ bool, _ bool, _ bool) error {
	return nil
}

type TriggerStub struct {
}

func (t TriggerStub) RunAccessChangeEvent() error {
	return nil
}

type ClientStub struct {
	response string
	status   int
}

func (c *ClientStub) Post(_, _ string, _ interface{}) (*http.Response, error) {
	if c.status != 200 {
		return nil, fmt.Errorf("error code: %v", c.status)
	}

	r := ioutil.NopCloser(bytes.NewReader([]byte(c.response)))
	return &http.Response{
		StatusCode: c.status,
		Body:       r,
	}, nil
}

type NetworkInfoStub struct {
}

func (n NetworkInfoStub) IPv6() *string {
	ip := "[::1]"
	return &ip
}

func TestOk_GoodResponse(t *testing.T) {
	client := &ClientStub{`{"success":true,"data":"OK"}`, 200}
	address := New(&UserConfigStub{}, &RedirectStub{}, &TriggerStub{}, client, &NetworkInfoStub{}, log.Default())
	ip := "1.1.1.1"
	err := address.Probe(&ip, 1)
	assert.Nil(t, err)
}

func TestFail_BadResponse(t *testing.T) {
	client := &ClientStub{`{"success":false,"message":"error"}`, 200}
	address := New(&UserConfigStub{}, &RedirectStub{}, &TriggerStub{}, client, &NetworkInfoStub{}, log.Default())
	ip := "1.1.1.1"
	err := address.Probe(&ip, 1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Unable to verify")

}

func TestFail_NotAPublicIp(t *testing.T) {
	client := &ClientStub{`{"success":false,"message":"error"}`, 200}
	address := New(&UserConfigStub{}, &RedirectStub{}, &TriggerStub{}, client, &NetworkInfoStub{}, log.Default())
	ip := "192.168.1.1"
	err := address.Probe(&ip, 1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "IP: 192.168.1.1 is not public")

}
