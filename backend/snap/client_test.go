package snap

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HttpClientStub struct {
	response string
	status   int
	err      error
	request  *http.Request
}

func (c *HttpClientStub) Get(_ string) (*http.Response, error) {
	r := io.NopCloser(bytes.NewReader([]byte(c.response)))
	return &http.Response{
		StatusCode: c.status,
		Body:       r,
	}, c.err
}

func (c *HttpClientStub) Post(_, _ string, _ io.Reader) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (c *HttpClientStub) Do(request *http.Request) (*http.Response, error) {
	c.request = request
	r := io.NopCloser(bytes.NewReader([]byte(c.response)))
	return &http.Response{
		StatusCode: c.status,
		Body:       r,
	}, c.err
}

func TestSnapdHttpClient_Get_404_IsError(t *testing.T) {
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
	snapd := &SnapdHttpClient{
		client: &HttpClientStub{response: json, status: 404},
	}
	_, err := snapd.Get("url")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, NotFound)
}

func TestSnapdHttpClient_Get_500_IsNotError(t *testing.T) {
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
	snapd := &SnapdHttpClient{
		client: &HttpClientStub{response: json, status: 500},
	}
	resp, err := snapd.Get("url")
	assert.Nil(t, err)
	assert.Equal(t, json, string(resp))
}

func TestSnapdHttpClient_Put_UsesPutMethod(t *testing.T) {
	client := &HttpClientStub{response: `{"status": "Accepted"}`, status: 202}
	snapd := &SnapdHttpClient{client: client}

	resp, err := snapd.Put("http://unix/v2/snaps/system/conf", "application/json", bytes.NewBufferString("{}"))

	assert.Nil(t, err)
	assert.Equal(t, 202, resp.StatusCode)
	assert.Equal(t, http.MethodPut, client.request.Method)
	assert.Equal(t, "application/json", client.request.Header.Get("Content-Type"))
}
