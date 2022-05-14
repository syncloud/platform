package snap

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
	"io"
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

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type DeviceInfo interface {
	Url(app string) string
}

type Config interface {
	Channel() string
}

type Snapd struct {
	client     SnapdClient
	deviceInfo DeviceInfo
	config     Config
	httpClient HttpClient
	logger     *zap.Logger
}

type Response struct {
	Result []model.Snap `json:"result"`
	Status string       `json:"status"`
}

func New(client SnapdClient, deviceInfo DeviceInfo, config Config, httpClient HttpClient, logger *zap.Logger) *Snapd {
	return &Snapd{
		client:     client,
		deviceInfo: deviceInfo,
		config:     config,
		httpClient: httpClient,
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
	bodyBytes, err := s.request("http://unix/v2/snaps")
	if err != nil {
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

func (s *Snapd) request(url string) ([]byte, error) {
	resp, err := s.client.Get(url)
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
	return bodyBytes, nil

}

func (s *Snapd) StoreSnaps() ([]model.Snap, error) {
	return s.Find("*")
}

func (s *Snapd) Installer() (*model.InstallerInfo, error) {
	s.logger.Info("installer")
	channel := s.config.Channel()
	systemInfoBytes, err := s.request(fmt.Sprintf("http://unix/v2/system-info"))
	if err != nil {
		return nil, err
	}
	var systemInfo model.SystemInfo
	err = json.Unmarshal(systemInfoBytes, &systemInfo)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}

	resp, err := s.httpClient.Get(fmt.Sprintf("http://apps.syncloud.org/releases/%s/snapd.version", channel))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &model.InstallerInfo{
		StoreVersion:     systemInfo.Result.Version,
		InstalledVersion: string(body),
	}, nil
}

func (s *Snapd) Find(query string) ([]model.Snap, error) {
	s.logger.Info("available snaps", zap.String("query", query))
	bodyBytes, err := s.request(fmt.Sprintf("http://unix/v2/find?name=%s", query))
	if err != nil {
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
