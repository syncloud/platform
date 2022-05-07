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

func TestProbe_Ok_GoodResponse(t *testing.T) {
	client := &ClientStub{`{"success":true,"data":"OK"}`, 200}
	probe := NewProbe(&UserConfigStub{}, client, log.Default())
	err := probe.Probe("1.1.1.1", 1)
	assert.Nil(t, err)
}

func TestProbe_Fail_BadResponse(t *testing.T) {
	client := &ClientStub{`{"success":false,"message":"error"}`, 200}
	probe := NewProbe(&UserConfigStub{}, client, log.Default())
	err := probe.Probe("1.1.1.1", 1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Unable to verify")
}

func TestProbe_Fail_NotAPublicIp(t *testing.T) {
	client := &ClientStub{`{"success":false,"message":"error"}`, 200}
	probe := NewProbe(&UserConfigStub{}, client, log.Default())
	err := probe.Probe("192.168.1.1", 1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "IP: 192.168.1.1 is not public")
}

func TestProbe_Fail_NotValidIp(t *testing.T) {
	client := &ClientStub{`{"success":false,"message":"error"}`, 200}
	probe := NewProbe(&UserConfigStub{}, client, log.Default())
	err := probe.Probe("1.1.1", 1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "IP: 1.1.1 is not valid")
}
