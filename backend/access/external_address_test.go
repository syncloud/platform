package access

import (
	"bytes"
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

func (r RedirectStub) Update(ipv4 *string, ipv6 *string, port *int, ipv4Enabled bool, ipv4Public bool, ipv6Enabled bool) error {
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

func (c *ClientStub) Post(_, _ string, body interface{}) (*http.Response, error) {
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

func TestOk(t *testing.T) {
	client := &ClientStub{`{"success":true,"data":"OK"}`, 200}
	address := New(&UserConfigStub{}, &RedirectStub{}, &TriggerStub{}, client, &NetworkInfoStub{}, log.Default())
	ip := "1.1.1.1"
	err := address.Probe(&ip, 1)
	assert.Nil(t, err)
}

func TestFail(t *testing.T) {
	client := &ClientStub{`{"success":false,"message":"error"}`, 500}
	address := New(&UserConfigStub{}, &RedirectStub{}, &TriggerStub{}, client, &NetworkInfoStub{}, log.Default())
	ip := "1.1.1.1"
	err := address.Probe(&ip, 1)
	assert.NotNil(t, err)

}
