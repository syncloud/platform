package snap

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
)

var NotFound = errors.New("app not found")

type SnapdHttpClient struct {
	client HttpClient
	logger *zap.Logger
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, bodyType string, body io.Reader) (*http.Response, error)
}

func NewSnapdHttpClient(logger *zap.Logger) *SnapdHttpClient {
	return &SnapdHttpClient{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", SOCKET)
				},
			},
		},
		logger: logger,
	}
}

func (c *SnapdHttpClient) Get(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		c.logger.Error("cannot connect", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFound
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("cannot read output", zap.Error(err))
		return nil, err
	}
	return bodyBytes, nil
}

func (c *SnapdHttpClient) Post(url, bodyType string, body io.Reader) (*http.Response, error) {
	return c.client.Post(url, bodyType, body)
}
