package snap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	SOCKET = "/var/run/snapd.socket"
)

type SnapdClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Snapd struct {
	client SnapdClient
}

type Response struct {
	Result []Snap `json:"result"`
}

func New(client SnapdClient) *Snapd {
	return &Snapd{
		client: client,
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
		return nil, err
	}

	var response Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return response.Result, nil

}
