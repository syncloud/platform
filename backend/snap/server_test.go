package snap

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

type ClientStub struct {
	json   string
	error  bool
	status int
}

func (c *ClientStub) Get(_ string) (*http.Response, error) {
	if c.error {
		return nil, fmt.Errorf("error")
	}
	r := io.NopCloser(bytes.NewReader([]byte(c.json)))
	return &http.Response{
		StatusCode: c.status,
		Body:       r,
	}, nil
}

type DeviceInfoStub struct {
}

func (d DeviceInfoStub) Url(app string) string {
	return fmt.Sprintf("%s.domain.tld", app)
}

type HttpClientStub struct {
	response string
	status   int
}

func (h HttpClientStub) Get(_ string) (*http.Response, error) {
	if h.status != 200 {
		return nil, fmt.Errorf("error code: %v", h.status)
	}

	r := io.NopCloser(bytes.NewReader([]byte(h.response)))
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
			"summary": "test summary",
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

	snapd := NewServer(&ClientStub{json: json, error: false, status: 200}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.Snaps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 1)
	assert.Equal(t, apps[0].Apps[0].Name, "test")
}

func TestInstalledSnaps_Error(t *testing.T) {

	snapd := NewServer(&ClientStub{error: true}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.Snaps()

	assert.Nil(t, apps)
	assert.NotNil(t, err)
}

func TestStoreSnaps_OK(t *testing.T) {
	json := `
{ 
	"result": [ 
		{ 
			"name": "test",
			"summary": "test summary",
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

	snapd := NewServer(&ClientStub{json: json, error: false, status: 200}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
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

	snapd := NewServer(&ClientStub{json: installed, error: false, status: 200}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{response: store, status: 200}, log.Default())
	installer, err := snapd.Installer()

	assert.Nil(t, err)
	assert.Equal(t, installer.InstalledVersion, "1")
	assert.Equal(t, installer.StoreVersion, "2")
}

func TestInstalledUserApps_OK(t *testing.T) {
	json := `
{ 
	"result": [ 
		{ 
			"name": "app",
			"summary": "app summary",
			"channel": "stable",
			"version": "1",
			"type": "app",
			"apps": [ 
				{
					"name": "app",
					"snap": "app"
				}
			]
		}, 
		{ 
			"name": "platform",
			"summary": "platform summary",
			"channel": "stable",
			"version": "1",
			"type": "system",
			"apps": [ 
				{
					"name": "platform",
					"snap": "platform"
				}
			]
		} 
	]
}
`

	snapd := NewServer(&ClientStub{json: json, error: false, status: 200}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.InstalledUserApps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 1)
	assert.Equal(t, apps[0].Id, "app")
}

func TestStoreUserApps_OK(t *testing.T) {
	json := `
{ 
	"result": [ 
		{ 
			"name": "app",
			"summary": "app summary",
			"channel": "stable",
			"version": "1",
			"type": "app",
			"apps": [ 
				{
					"name": "app",
					"snap": "app"
				}
			]
		}, 
		{ 
			"name": "platform",
			"summary": "platform summary",
			"channel": "stable",
			"version": "1",
			"type": "system",
			"apps": [ 
				{
					"name": "platform",
					"snap": "platform"
				}
			]
		}  
	]
}
`

	snapd := NewServer(&ClientStub{json: json, error: false, status: 200}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.StoreUserApps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 1)
	assert.Equal(t, apps[0].Id, "app")
}

func TestServer_FindInStore_Found(t *testing.T) {
	json := `
{ 
	"status": "OK",
	"result": [ 
		{ 
			"name": "app",
			"summary": "app summary",
			"channel": "stable",
			"version": "1",
			"type": "app",
			"apps": [ 
				{
					"name": "app",
					"snap": "app"
				}
			]
		}
	]
}
`

	client := &ClientStub{json: json, error: false, status: 200}
	snapd := NewServer(client, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.FindInStore("app")

	assert.Nil(t, err)
	assert.Equal(t, "app", found.App.Id)
}

func TestServer_FindInStore_NotFound(t *testing.T) {
	json := `
{ 
	"status": "Error",
	"result": {
		"message": "not found"
	}
}
`

	client := &ClientStub{json: json, error: false, status: 500}
	snapd := NewServer(client, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.FindInStore("app")

	assert.Nil(t, err)
	assert.Nil(t, found)
}

func TestServer_Changes_Error(t *testing.T) {
	json := `
{
    "type": "error",
    "status-code": 401,
    "status": "Unauthorized",
    "result": {
        "message": "access denied",
        "kind": "login-required",
    }
}
`

	snapd := NewServer(&ClientStub{json: json, error: false, status: 200}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	_, err := snapd.Changes()

	assert.NotNil(t, err)
}

func TestServer_Changes_True(t *testing.T) {
	json := `
{
    "type": "sync",
    "status-code": 200,
    "status": "OK",
    "result": [
		{
			"id": "123"
		}
	]
}
`

	snapd := NewServer(&ClientStub{json: json, error: false, status: 200}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	progress, err := snapd.Changes()

	assert.Nil(t, err)
	assert.True(t, progress.IsRunning)
}

func TestServer_Changes_False(t *testing.T) {
	json := `
{
    "type": "sync",
    "status-code": 200,
    "status": "OK",
    "result": []
}
`

	snapd := NewServer(&ClientStub{json: json, error: false, status: 200}, &DeviceInfoStub{}, &ConfigStub{}, &HttpClientStub{}, log.Default())
	progress, err := snapd.Changes()

	assert.Nil(t, err)
	assert.False(t, progress.IsRunning)
}
