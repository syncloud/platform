package snap

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"io/ioutil"
	"net/http"
	"testing"
)

type ClientStub struct {
	error bool
}

func (c *ClientStub) Get(_ string) (resp *http.Response, err error) {
	if c.error {
		return nil, fmt.Errorf("error")
	}
	json := `
{ 
	"result": [ 
		{ 
			"name": "test",
			"summary": "test",
			"channel": "stable",
			"version": "1",
			"apps": [ 
				{
					"name": "test",
					"snap": "test"
				}
			]
		} 
	]
}
`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

type DeviceInfoStub struct {
}

func (d DeviceInfoStub) Url(_ string) string {
	//TODO implement me
	panic("implement me")
}

type HttpClientStub struct {
	response string
	status   int
}

func (h HttpClientStub) Post(url, bodyType string, body interface{}) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (h HttpClientStub) Get(url string) (*http.Response, error) {
	if h.status != 200 {
		return nil, fmt.Errorf("error code: %v", h.status)
	}

	r := ioutil.NopCloser(bytes.NewReader([]byte(h.response)))
	return &http.Response{
		StatusCode: h.status,
		Body:       r,
	}, nil
}

type ConfigStub struct {
}

func (c ConfigStub) Channel() string {
	return "stable"
}

func TestAppsOK(t *testing.T) {

	snapd := New(&ClientStub{error: false}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.InstalledSnaps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 1)
	assert.Equal(t, apps[0].Apps[0].Name, "test")
}

func TestAppsError(t *testing.T) {

	snapd := New(&ClientStub{error: true}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.InstalledSnaps()

	assert.Nil(t, apps)
	assert.NotNil(t, err)
}
