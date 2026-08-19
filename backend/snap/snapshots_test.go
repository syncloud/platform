package snap

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

type SnapshotsClientStub struct {
	listJson  string
	confJson  string
	getError  error
	posted    []string
	putBodies []string
	status    int
}

func (c *SnapshotsClientStub) Get(url string) ([]byte, error) {
	if c.getError != nil {
		return nil, c.getError
	}
	if strings.HasPrefix(url, "http://unix/v2/snaps/system/conf") {
		return []byte(c.confJson), nil
	}
	return []byte(c.listJson), nil
}

func (c *SnapshotsClientStub) Post(_, _ string, body io.Reader) (*http.Response, error) {
	return c.record(&c.posted, body)
}

func (c *SnapshotsClientStub) Put(_, _ string, body io.Reader) (*http.Response, error) {
	return c.record(&c.putBodies, body)
}

func (c *SnapshotsClientStub) record(to *[]string, body io.Reader) (*http.Response, error) {
	content, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	*to = append(*to, string(content))
	status := c.status
	if status == 0 {
		status = http.StatusAccepted
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(`{"status": "Accepted"}`)),
	}, nil
}

const autoSets = `
{
  "result": [
    { "id": 42, "snapshots": [ { "set": 42, "snap": "jellyfin", "size": 186000, "auto": true } ] },
    { "id": 43, "snapshots": [ { "set": 43, "snap": "nextcloud", "size": 900, "auto": false } ] },
    { "id": 44, "snapshots": [ { "set": 44, "snap": "users", "size": 24923, "auto": true } ] }
  ]
}
`

const retentionNotSet = `
{
  "type": "error",
  "status-code": 400,
  "result": { "message": "snap \"core\" has no \"snapshots\" configuration option", "kind": "option-not-found" }
}
`

const retentionDisabled = `
{
  "type": "sync",
  "status-code": 200,
  "result": { "snapshots.automatic.retention": "no" }
}
`

func TestSnapshots_ForgetAutoOnly(t *testing.T) {
	client := &SnapshotsClientStub{listJson: autoSets}
	snapshots := NewSnapshots(client, log.Default())

	assert.NoError(t, snapshots.ForgetAuto())
	assert.Equal(t, []string{
		`{"action":"forget","set":42}`,
		`{"action":"forget","set":44}`,
	}, client.posted)
}

func TestSnapshots_ForgetNothingWhenEmpty(t *testing.T) {
	client := &SnapshotsClientStub{listJson: `{"result": []}`}
	snapshots := NewSnapshots(client, log.Default())

	assert.NoError(t, snapshots.ForgetAuto())
	assert.Empty(t, client.posted)
}

func TestSnapshots_ForgetFailsOnSnapdError(t *testing.T) {
	client := &SnapshotsClientStub{listJson: autoSets, status: http.StatusInternalServerError}
	snapshots := NewSnapshots(client, log.Default())

	assert.Error(t, snapshots.ForgetAuto())
}

func TestSnapshots_DisableAutomaticWhenNotSet(t *testing.T) {
	client := &SnapshotsClientStub{confJson: retentionNotSet}
	snapshots := NewSnapshots(client, log.Default())

	assert.NoError(t, snapshots.DisableAutomatic())
	assert.Equal(t, []string{`{"snapshots":{"automatic":{"retention":"no"}}}`}, client.putBodies)
}

func TestSnapshots_DisableAutomaticIsIdempotent(t *testing.T) {
	client := &SnapshotsClientStub{confJson: retentionDisabled}
	snapshots := NewSnapshots(client, log.Default())

	assert.NoError(t, snapshots.DisableAutomatic())
	assert.Empty(t, client.putBodies)
}
