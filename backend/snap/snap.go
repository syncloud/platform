package snap

import (
	"context"
	"net"
	"net/http"
)

const (
	SOCKET = "/var/run/snapd.socket"
)

type Snap struct {
	client http.Client
}

func New() *Snap {
	return &Snap{
		client: http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", SOCKET)
				},
			},
		},
	}
}

func (snap *Snap) ListAllApps() {
	snap.client.Get("http://unix/v2/snaps")
}
