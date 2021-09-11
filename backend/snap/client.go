package snap

import (
	"context"
	"net"
	"net/http"
)

func NewClient() SnapdClient {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", SOCKET)
			},
		},
	}
}
