package snap

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

type ClientStub struct {
	snapsJson   string
	snapsError  error
	snapJson    string
	snapError   error
	findJson    string
	findError   error
	systemJson  string
	systemError error
}

func (c *ClientStub) Post(_, _ string, _ io.Reader) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClientStub) Get(url string) ([]byte, error) {
	if strings.HasPrefix(url, "http://unix/v2/snaps/") {
		return []byte(c.snapJson), c.snapError
	}
	if strings.HasPrefix(url, "http://unix/v2/snaps") {
		return []byte(c.snapsJson), c.snapsError
	}
	if strings.HasPrefix(url, "http://unix/v2/find") {
		return []byte(c.findJson), c.findError
	}
	if strings.HasPrefix(url, "http://unix/v2/system-info") {
		return []byte(c.systemJson), c.systemError
	}
	panic("implement me")
}

type UserConfigStub struct {
}

func (d UserConfigStub) Url(app string) string {
	return fmt.Sprintf("%s.domain.tld", app)
}

type SystemConfigStub struct {
}

func (c SystemConfigStub) Channel() string {
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

	snapd := NewServer(&ClientStub{snapsJson: json}, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.Snaps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 1)
	assert.Equal(t, apps[0].Apps[0].Name, "test")
}

func TestInstalledSnaps_Error(t *testing.T) {

	snapd := NewServer(&ClientStub{snapsError: fmt.Errorf("error")}, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
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

	snapd := NewServer(&ClientStub{findJson: json}, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
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

	snapd := NewServer(&ClientStub{systemJson: installed}, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{response: store, status: 200}, log.Default())
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

	snapd := NewServer(&ClientStub{snapsJson: json}, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.InstalledUserApps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 1)
	assert.Equal(t, apps[0].Id, "app")
}

func TestInstalledUserApps_Sorted(t *testing.T) {
	json := `
{ 
	"result": [ 
		{ 
			"name": "app2",
			"summary": "app2 summary",
			"channel": "stable",
			"version": "1",
			"type": "app",
			"apps": [ 
				{
					"name": "app2",
					"snap": "app2"
				}
			]
		}, 
		{ 
			"name": "app1",
			"summary": "app1 summary",
			"channel": "stable",
			"version": "1",
			"type": "app",
			"apps": [ 
				{
					"name": "app1",
					"snap": "app1"
				}
			]
		} 
	]
}
`

	snapd := NewServer(&ClientStub{snapsJson: json}, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	apps, err := snapd.InstalledUserApps()

	assert.Nil(t, err)
	assert.Equal(t, len(apps), 2)
	assert.Equal(t, apps[0].Id, "app1")
	assert.Equal(t, apps[1].Id, "app2")
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

	snapd := NewServer(&ClientStub{findJson: json}, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
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

	client := &ClientStub{findJson: json}
	snapd := NewServer(client, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.FindInStore("app")

	assert.Nil(t, err)
	assert.Equal(t, "app", found.App.Id)
	assert.Equal(t, "1", *found.CurrentVersion)
	assert.Nil(t, found.InstalledVersion)
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

	client := &ClientStub{findJson: json}
	snapd := NewServer(client, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.FindInStore("app")

	assert.Nil(t, err)
	assert.Nil(t, found)
}

func TestServer_FindInstalled_Found(t *testing.T) {
	json := `
{
  "type": "sync",
  "status-code": 200,
  "status": "OK",
  "result": {
    "id": "mail.239",
    "summary": "Mail server",
    "description": "Mail",
    "name": "mail",
    "status": "active",
    "type": "app",
    "version": "239",
    "channel": "stable",
    "revision": "239",
    "private": false,
    "devmode": false,
    "jailmode": false,
    "apps": [],
    "contact": "",
    "install-date": "2022-08-24T23:45:26Z"
  }
}
`

	client := &ClientStub{snapJson: json}
	snapd := NewServer(client, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.FindInstalled("mail")

	assert.Nil(t, err)
	assert.Equal(t, "mail", found.Name)
}

func TestServer_FindInstalled_NotFound(t *testing.T) {
	json := `
{
	"type":"error",
	"status-code":404,
	"status":"Not Found",
	"result":{
		"message":"snap not installed",
		"kind":"snap-not-found",
		"value":"files"
	}
}
`

	client := &ClientStub{snapJson: json, snapError: NotFound}
	snapd := NewServer(client, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.FindInstalled("files")

	assert.Nil(t, err)
	assert.Nil(t, found)
}

func TestServer_Find_NotInstalled(t *testing.T) {
	snapJson := `
{
	"type":"error",
	"status-code":404,
	"status":"Not Found",
	"result":{
		"message":"snap not installed",
		"kind":"snap-not-found",
		"value":"files"
	}
}
`

	findJson := `
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
	client := &ClientStub{snapJson: snapJson, snapError: NotFound, findJson: findJson}
	snapd := NewServer(client, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.Find("app")

	assert.Nil(t, err)
	assert.Equal(t, "1", *found.CurrentVersion)
	assert.Nil(t, found.InstalledVersion)
}

func TestServer_Find_Installed(t *testing.T) {
	snapJson := `
{
  "type": "sync",
  "status-code": 200,
  "status": "OK",
  "result": {
    "id": "app.239",
    "summary": "app summary",
    "description": "",
    "name": "app",
    "status": "active",
    "type": "app",
    "version": "1",
    "channel": "stable",
    "revision": "239",
    "apps": [],
    "contact": "",
    "install-date": "2022-08-24T23:45:26Z"
  }
}
`

	findJson := `
{ 
	"status": "OK",
	"result": [ 
		{ 
			"name": "app",
			"summary": "app summary",
			"channel": "stable",
			"version": "2",
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
	client := &ClientStub{snapJson: snapJson, findJson: findJson}
	snapd := NewServer(client, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.Find("app")

	assert.Nil(t, err)
	assert.Equal(t, "2", *found.CurrentVersion)
	assert.Equal(t, "1", *found.InstalledVersion)
}

func TestServer_Find_NotInStore(t *testing.T) {
	snapJson := `
{
  "type": "sync",
  "status-code": 200,
  "status": "OK",
  "result": {
    "id": "app.239",
    "summary": "app summary",
    "description": "",
    "name": "app",
    "status": "active",
    "type": "app",
    "version": "1",
    "channel": "stable",
    "revision": "239",
    "apps": [],
    "contact": "",
    "install-date": "2022-08-24T23:45:26Z"
  }
}
`

	findJson := `
{ 
	"status": "Error",
	"result": {
		"message": "not found"
	}
}
`
	client := &ClientStub{snapJson: snapJson, findJson: findJson}
	snapd := NewServer(client, &SystemConfigStub{}, &UserConfigStub{}, &HttpClientStub{}, log.Default())
	found, err := snapd.Find("app")

	assert.Nil(t, err)
	assert.Nil(t, found.CurrentVersion)
	assert.Equal(t, "1", *found.InstalledVersion)
}
