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
	json  string
	error bool
}

func (c *ClientStub) Get(_ string) (resp *http.Response, err error) {
	if c.error {
		return nil, fmt.Errorf("error")
	}
	r := ioutil.NopCloser(bytes.NewReader([]byte(c.json)))
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

func TestInstalledSnaps_OK(t *testing.T) {
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

	snapd := New(&ClientStub{json: json, error: false}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.InstalledSnaps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 1)
	assert.Equal(t, apps[0].Apps[0].Name, "test")
}

func TestInstalledSnaps_Error(t *testing.T) {

	snapd := New(&ClientStub{error: true}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.InstalledSnaps()

	assert.Nil(t, apps)
	assert.NotNil(t, err)
}

func TestStoreSnaps_OK(t *testing.T) {
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

	snapd := New(&ClientStub{json: json, error: false}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.StoreSnaps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 1)
	assert.Equal(t, apps[0].Apps[0].Name, "test")
}

func TestInstaller_OK(t *testing.T) {
	installed := `
{ 
	"result": { 
		"version": "1"
	} 
}
`
	store := "2"

	snapd := New(&ClientStub{json: installed, error: false}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{response: store, status: 200}, log.Default())
	installer, err := snapd.Installer()

	assert.Nil(t, err)
	assert.Equal(t, installer.InstalledVersion, "1")
	assert.Equal(t, installer.StoreVersion, "2")
}
