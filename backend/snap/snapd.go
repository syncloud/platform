package snap

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sort"
)

const (
	SOCKET = "/var/run/snapd.socket"
)

type SnapdClient interface {
	Get(url string) (resp *http.Response, err error)
}

type DeviceInfo interface {
	Url(app string) string
}

type Snapd struct {
	client     SnapdClient
	deviceInfo DeviceInfo
	logger     *zap.Logger
}

type Response struct {
	Result []model.Snap `json:"result"`
	Status string       `json:"status"`
}

func New(client SnapdClient, deviceInfo DeviceInfo, logger *zap.Logger) *Snapd {
	return &Snapd{
		client:     client,
		deviceInfo: deviceInfo,
		logger:     logger,
	}
}

func (s *Snapd) InstalledUserApps() ([]model.SyncloudApp, error) {
	snaps, err := s.InstalledSnaps()
	if err != nil {
		return nil, err
	}
	var apps []model.SyncloudApp
	for _, snap := range snaps {
		if snap.IsApp() {
			app := snap.ToInstalledApp(s.deviceInfo.Url(snap.Name))
			apps = append(apps, app.App)
		}
	}
	return apps, nil
}

func (s *Snapd) StoreUserApps() ([]model.SyncloudApp, error) {
	snaps, err := s.StoreSnaps()
	if err != nil {
		return nil, err
	}
	var apps []model.SyncloudApp
	for _, snap := range snaps {
		if snap.IsApp() {
			app := snap.ToStoreApp(s.deviceInfo.Url(snap.Name))
			apps = append(apps, app.App)
		}
	}
	return apps, nil
}

func (s *Snapd) InstalledSnaps() ([]model.Snap, error) {
	resp, err := s.client.Get("http://unix/v2/snaps")
	if err != nil {
		s.logger.Error("cannot connect", zap.Error(err))
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("status", zap.Error(err))
		return nil, fmt.Errorf("unable to get apps list, status code: %d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("cannot read output", zap.Error(err))
		return nil, err
	}

	var response Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}
	return response.Result, nil

}
func (s *Snapd) StoreSnaps() ([]model.Snap, error) {
	return s.Find("*")
}

func (s *Snapd) Find(query string) ([]model.Snap, error) {

	s.logger.Info("available snaps", zap.String("query", query))
	resp, err := s.client.Get(fmt.Sprintf("http://unix/v2/find?name=%s", query))
	if err != nil {
		s.logger.Error("cannot connect", zap.Error(err))
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("status", zap.Error(err))
		return nil, fmt.Errorf("unable to get apps list, status code: %d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("cannot read output", zap.Error(err))
		return nil, err
	}
	var response Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}

	if query != "*" && response.Status != "OK" {
		return make([]model.Snap, 0), nil
	}

	sort.SliceStable(response.Result, func(i, j int) bool {
		return response.Result[i].Name < response.Result[i].Name
	})
	return response.Result, nil
}
