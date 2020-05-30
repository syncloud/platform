package snapd

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

type Snapd struct {
	client http.Client
}

type Response struct {
	Result []Snap `json:"result"`
}

func New() *Snapd {
	return &Snapd{
		client: http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", SOCKET)
				},
			},
		},
	}
}

func (snap *Snapd) ListAllApps() ([]Snap, error) {
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
	var response Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return response.Result, nil

}
