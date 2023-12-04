package snap

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

type ChangesHttpClientStub struct {
	json  string
	error bool
}

func (c *ChangesHttpClientStub) Post(_, _ string, _ io.Reader) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChangesHttpClientStub) Get(url string) ([]byte, error) {
	if c.error {
		return nil, fmt.Errorf("error")
	}
	return []byte(c.json), nil
}

func TestChangesClient_Changes_Error(t *testing.T) {
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

	snapd := NewChangesClient(&ChangesHttpClientStub{json: json}, log.Default())
	_, err := snapd.Changes()

	assert.NotNil(t, err)
}

func TestChangesClient_Changes_True(t *testing.T) {
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

	snapd := NewChangesClient(&ChangesHttpClientStub{json: json}, log.Default())
	progress, err := snapd.Changes()

	assert.Nil(t, err)
	assert.True(t, progress.IsRunning)
}

func TestChangesClient_Changes_False(t *testing.T) {
	json := `
{
    "type": "sync",
    "status-code": 200,
    "status": "OK",
    "result": []
}
`

	snapd := NewChangesClient(&ChangesHttpClientStub{json: json}, log.Default())
	progress, err := snapd.Changes()

	assert.Nil(t, err)
	assert.False(t, progress.IsRunning)
}

func TestChangesClient_Changes_Progress(t *testing.T) {
	json := `
{
  "type": "sync",
  "status-code": 200,
  "status": "OK",
  "result": [
    {
      "id": "2161",
      "kind": "refresh-snap",
      "summary": "Refresh \"matrix\" snap from \"latest/stable\" channel",
      "status": "Doing",
      "tasks": [
        {
          "id": "11282",
          "kind": "download-snap",
          "summary": "Download snap \"matrix\" (281) from channel \"latest/stable\"",
          "status": "Doing",
          "progress": {
            "label": "matrix",
            "done": 146936354,
            "total": 312000512
          },
          "spawn-time": "2023-11-30T16:53:24.663070859Z"
        },
        {
          "id": "11300",
          "kind": "run-hook",
          "summary": "Run configure hook of \"matrix\" snap if present",
          "status": "Do",
          "progress": {
            "label": "",
            "done": 0,
            "total": 1
          },
          "spawn-time": "2023-11-30T16:53:24.663804275Z"
        }
      ],
      "ready": false,
      "spawn-time": "2023-11-30T16:53:24.663995517Z"
    }
  ]
}
`

	snapd := NewChangesClient(&ChangesHttpClientStub{json: json}, log.Default())
	progress, err := snapd.Changes()

	assert.Nil(t, err)
	assert.True(t, progress.IsRunning)
	assert.Equal(t, false, progress.Progress["matrix"].Indeterminate)
	assert.Equal(t, int64(47), progress.Progress["matrix"].Percentage)
	assert.Equal(t, "Downloading", progress.Progress["matrix"].Summary)
}

func TestChangesClient_Changes_Progress_Unknown_Zero(t *testing.T) {
	json := `
{
  "type": "sync",
  "status-code": 200,
  "status": "OK",
  "result": [
    {
      "id": "2161",
      "kind": "refresh-snap",
      "summary": "Refresh \"matrix\" snap from \"latest/stable\" channel",
      "status": "Doing",
      "tasks": [
        {
          "id": "11282",
          "kind": "download-snap",
          "summary": "Download snap \"matrix\" (281) from channel \"latest/stable\"",
          "status": "Doing",
          "progress": {
            "label": "matrix",
            "done": 1,
            "total": 1
          },
          "spawn-time": "2023-11-30T16:53:24.663070859Z"
        },
        {
          "id": "11300",
          "kind": "run-hook",
          "summary": "Run configure hook of \"matrix\" snap if present",
          "status": "Do",
          "progress": {
            "label": "",
            "done": 0,
            "total": 1
          },
          "spawn-time": "2023-11-30T16:53:24.663804275Z"
        }
      ],
      "ready": false,
      "spawn-time": "2023-11-30T16:53:24.663995517Z"
    }
  ]
}
`

	snapd := NewChangesClient(&ChangesHttpClientStub{json: json}, log.Default())
	progress, err := snapd.Changes()

	assert.Nil(t, err)
	assert.True(t, progress.IsRunning)
	assert.Equal(t, false, progress.Progress["matrix"].Indeterminate)
	assert.Equal(t, int64(0), progress.Progress["matrix"].Percentage)
	assert.Equal(t, "Downloading", progress.Progress["matrix"].Summary)
}

func TestChangesClient_Changes_Indeterminate(t *testing.T) {
	json := `
{
  "type": "sync",
  "status-code": 200,
  "status": "OK",
  "result": [
    {
      "id": "2161",
      "kind": "refresh-snap",
      "summary": "Refresh \"matrix\" snap from \"latest/stable\" channel",
      "status": "Doing",
      "tasks": [
		{
          "id": "11282",
          "kind": "download-snap",
          "summary": "Download snap \"matrix\" (281) from channel \"latest/stable\"",
          "status": "Done",
          "progress": {
            "label": "matrix",
            "done": 312000512,
            "total": 312000512
          },
          "spawn-time": "2023-11-30T16:53:24.663070859Z",
          "ready-time": "2023-11-30T16:55:31.819476646Z"
        },
        {
          "id": "11285",
          "kind": "run-hook",
          "summary": "Run pre-refresh hook of \"matrix\" snap if present",
          "status": "Doing",
          "progress": {
            "label": "",
            "done": 1,
            "total": 1
          },
          "spawn-time": "2023-11-30T16:53:24.663159158Z"
        }
      ],
      "ready": false,
      "spawn-time": "2023-11-30T16:53:24.663995517Z"
    }
  ]
}
`

	snapd := NewChangesClient(&ChangesHttpClientStub{json: json}, log.Default())
	progress, err := snapd.Changes()

	assert.Nil(t, err)
	assert.True(t, progress.IsRunning)
	assert.Equal(t, true, progress.Progress["matrix"].Indeterminate)
	assert.Equal(t, int64(0), progress.Progress["matrix"].Percentage)
	assert.Equal(t, "Upgrading", progress.Progress["matrix"].Summary)
}

func TestChangesClient_Changes_Simultaneous(t *testing.T) {
	json := `
{
  "type": "sync",
  "status-code": 200,
  "status": "OK",
  "result": [
    {
      "id": "487",
      "kind": "install-snap",
      "summary": "Install \"jellyfin\" snap",
      "status": "Doing",
      "tasks": [
        {
          "id": "2887",
          "kind": "download-snap",
          "summary": "Download snap \"jellyfin\" (161) from channel \"stable\"",
          "status": "Doing",
          "progress": {
            "label": "jellyfin",
            "done": 105143265,
            "total": 192454656
          },
          "spawn-time": "2023-12-04T16:26:47.314203858Z"
        },
        {
          "id": "2896",
          "kind": "run-hook",
          "summary": "Run install hook of \"jellyfin\" snap if present",
          "status": "Do",
          "progress": {
            "label": "",
            "done": 0,
            "total": 1
          },
          "spawn-time": "2023-12-04T16:26:47.314762264Z"
        }
      ],
      "ready": false,
      "spawn-time": "2023-12-04T16:26:47.314887646Z"
    },
    {
      "id": "488",
      "kind": "install-snap",
      "summary": "Install \"collabora\" snap",
      "status": "Doing",
      "tasks": [
        {
          "id": "2901",
          "kind": "download-snap",
          "summary": "Download snap \"collabora\" (36) from channel \"stable\"",
          "status": "Done",
          "progress": {
            "label": "collabora",
            "done": 420900864,
            "total": 420900864
          },
          "spawn-time": "2023-12-04T16:26:55.201215861Z"
        },        
        {
          "id": "2910",
          "kind": "run-hook",
          "summary": "Run install hook of \"collabora\" snap if present",
          "status": "Doing",
          "progress": {
            "label": "",
            "done": 0,
            "total": 1
          },
          "spawn-time": "2023-12-04T16:26:55.201769017Z"
        }
      ],
      "ready": false,
      "spawn-time": "2023-12-04T16:26:55.201920442Z"
    }
  ]
}
`

	snapd := NewChangesClient(&ChangesHttpClientStub{json: json}, log.Default())
	progress, err := snapd.Changes()

	assert.Nil(t, err)
	assert.True(t, progress.IsRunning)

	assert.Equal(t, false, progress.Progress["jellyfin"].Indeterminate)
	assert.Equal(t, int64(54), progress.Progress["jellyfin"].Percentage)
	assert.Equal(t, "Downloading", progress.Progress["jellyfin"].Summary)

	assert.Equal(t, true, progress.Progress["collabora"].Indeterminate)
	assert.Equal(t, int64(0), progress.Progress["collabora"].Percentage)
	assert.Equal(t, "Installing", progress.Progress["collabora"].Summary)

}
