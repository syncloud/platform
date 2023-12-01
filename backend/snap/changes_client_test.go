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
	assert.Equal(t, false, progress.Progress.Indeterminate)
	assert.Equal(t, int64(47), progress.Progress.Percentage)
	assert.Equal(t, "Downloading", progress.Progress.Summary)
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
	assert.Equal(t, true, progress.Progress.Indeterminate)
	assert.Equal(t, int64(50), progress.Progress.Percentage)
	assert.Equal(t, "Upgrading", progress.Progress.Summary)
}
