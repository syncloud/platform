package snap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

const (
	SOCKET = "/var/run/snapd.socket"
)

type Snap struct {
	client http.Client
}

type Apps struct {
	Result []App `json:"result"`
}

type App struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Channel string `json:"channel"`
	Version string `json:"version"`
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

func (snap *Snap) ListAllApps() ([]App, error) {
	resp, err := snap.client.Get("http://unix/v2/snaps")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return nil, fmt.Errorf("unable to get apps list, status code: %d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var apps Apps
	err = json.Unmarshal(bodyBytes, &apps)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return apps.Result, nil

}
