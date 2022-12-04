package snap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
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

type Server struct {
	client     SnapdClient
	deviceInfo DeviceInfo
	config     Config
	httpClient HttpClient
	logger     *zap.Logger
}

func NewServer(client SnapdClient, deviceInfo DeviceInfo, config Config, httpClient HttpClient, logger *zap.Logger) *Server {
	return &Server{
		client:     client,
		deviceInfo: deviceInfo,
		config:     config,
		httpClient: httpClient,
		logger:     logger,
	}
}

func (s *Server) InstalledUserApps() ([]model.SyncloudApp, error) {
	snaps, err := s.Snaps()
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

func (s *Server) StoreUserApps() ([]model.SyncloudApp, error) {
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

func (s *Server) Snaps() ([]model.Snap, error) {
	bodyBytes, err := s.request("http://unix/v2/snaps")
	if err != nil {
		return nil, err
	}
	var response model.SnapsResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}
	return response.Result, nil

}

func (s *Server) Snap(name string) (model.Snap, error) {
	bodyBytes, err := s.request(fmt.Sprintf("http://unix/v2/snaps/%s", name))
	if err != nil {
		return model.Snap{}, err
	}
	var response model.SnapResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return model.Snap{}, err
	}
	return response.Result, nil

}

func (s *Server) request(url string) ([]byte, error) {
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
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("cannot read output", zap.Error(err))
		return nil, err
	}
	return bodyBytes, nil

}

func (s *Server) StoreSnaps() ([]model.Snap, error) {
	return s.find("*")
}

func (s *Server) Installer() (*model.InstallerInfo, error) {
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
		StoreVersion:     string(body),
		InstalledVersion: systemInfo.Result.Version,
	}, nil
}

func (s *Server) find(query string) ([]model.Snap, error) {
	s.logger.Info("available snaps", zap.String("query", query))
	bodyBytes, err := s.request(fmt.Sprintf("http://unix/v2/find?name=%s", query))
	if err != nil {
		return nil, err
	}
	var response model.SnapsResponse
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

func (s *Server) Changes() (bool, error) {
	s.logger.Info("snap changes")

	bodyBytes, err := s.request("http://unix/v2/changes?select=in-progress")
	if err != nil {
		return false, err
	}
	var response model.ServerResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return false, err
	}
	if response.Status != "OK" {
		var errorResponse model.ServerError
		err = json.Unmarshal(response.Result, &errorResponse)
		if err != nil {
			s.logger.Error("cannot unmarshal", zap.Error(err))
			return false, err
		}

		return false, fmt.Errorf(errorResponse.Message)
	}

	var changesResponse []model.Change
	err = json.Unmarshal(response.Result, &changesResponse)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return false, err
	}

	return len(changesResponse) > 0, nil
}
