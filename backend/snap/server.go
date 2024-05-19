package snap

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Get(url string) ([]byte, error)
	Post(url, bodyType string, body io.Reader) (*http.Response, error)
}

type DeviceInfo interface {
	Url(app string) string
}

type SystemConfig interface {
	Channel() string
}

type UserConfig interface {
	Url(app string) string
}

type ExternalHttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Server struct {
	client       SnapdClient
	systemConfig SystemConfig
	userConfig   UserConfig
	httpClient   ExternalHttpClient
	logger       *zap.Logger
}

func NewServer(
	client SnapdClient,
	systemConfig SystemConfig,
	userConfig UserConfig,
	httpClient ExternalHttpClient,
	logger *zap.Logger,
) *Server {
	return &Server{
		client:       client,
		systemConfig: systemConfig,
		userConfig:   userConfig,
		httpClient:   httpClient,
		logger:       logger,
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
			app := snap.ToInstalledApp(s.userConfig.Url(snap.Name))
			apps = append(apps, app.App)
		}
	}
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Id < apps[j].Id
	})
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
			app := snap.ToStoreApp(s.userConfig.Url(snap.Name))
			apps = append(apps, app.App)
		}
	}
	return apps, nil
}

func (s *Server) Snaps() ([]model.Snap, error) {
	bodyBytes, err := s.client.Get("http://unix/v2/snaps")
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

func (s *Server) FindInstalled(name string) (*model.Snap, error) {
	bodyBytes, err := s.client.Get(fmt.Sprintf("http://unix/v2/snaps/%s", name))
	if err != nil {
		if errors.Is(err, NotFound) {
			return nil, nil
		}
		return nil, err
	}

	var response model.SnapResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return &response.Result, nil

}

func (s *Server) Upgrade(name string) error {
	response, err := s.snapsAction("refresh", name)
	if err != nil {
		return err
	}
	if response.Status != "Accepted" {
		var serverError model.ServerError
		err = json.Unmarshal(response.Result, &serverError)
		if err != nil {
			return err
		}
		return fmt.Errorf(serverError.Message)
	}
	return nil
}

func (s *Server) Install(name string) error {
	_, err := s.snapsAction("install", name)
	return err
}

func (s *Server) Remove(name string) error {
	_, err := s.snapsAction("remove", name)
	return err
}

func (s *Server) snapsAction(action, name string) (*model.ServerResponse, error) {
	requestJson, err := json.Marshal(model.InstallRequest{Action: action})
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(fmt.Sprintf("http://unix/v2/snaps/%s", name), "application/json", bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFound
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("cannot read output", zap.Error(err))
		return nil, err
	}
	var response model.ServerResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}

	return &response, nil
}

func (s *Server) StoreSnaps() ([]model.Snap, error) {
	return s.find("*")
}

func (s *Server) Installer() (*model.InstallerInfo, error) {
	s.logger.Info("installer")
	channel := s.systemConfig.Channel()
	systemInfoBytes, err := s.client.Get(fmt.Sprintf("http://unix/v2/system-info"))
	if err != nil {
		return nil, err
	}
	var systemInfo model.SystemInfo
	err = json.Unmarshal(systemInfoBytes, &systemInfo)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}

	resp, err := s.httpClient.Get(fmt.Sprintf("http://apps.syncloud.org/releases/%s/snapd2.version", channel))
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
	s.logger.Info("find", zap.String("query", query))
	bodyBytes, err := s.client.Get(fmt.Sprintf("http://unix/v2/find?name=%s", query))
	if err != nil {
		return nil, err
	}
	var response model.ServerResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}

	if query != "*" && response.Status != "OK" {
		return make([]model.Snap, 0), nil
	}

	var snaps []model.Snap
	err = json.Unmarshal(response.Result, &snaps)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err), zap.String("response", string(bodyBytes)))
		return nil, err
	}
	return snaps, nil
}

func (s *Server) FindInStore(name string) (*model.SyncloudAppVersions, error) {
	s.logger.Info("find in store")
	found, err := s.find(name)
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		s.logger.Warn("No app found")
		return nil, nil
	}

	if len(found) > 1 {
		s.logger.Warn("More than one app found")
	}
	snap := found[0]
	installedApp := snap.ToStoreApp(s.userConfig.Url(snap.Name))
	return &installedApp, nil
}

func (s *Server) Find(name string) (*model.SyncloudAppVersions, error) {
	foundInstalledApp, err := s.FindInstalled(name)
	if err != nil {
		return nil, err
	}
	storeApp, err := s.FindInStore(name)
	if err != nil {
		return nil, err
	}
	if foundInstalledApp == nil && storeApp == nil {
		return nil, fmt.Errorf("not found")
	}

	if foundInstalledApp == nil {
		return storeApp, nil
	}

	installedApp := foundInstalledApp.ToInstalledApp(s.userConfig.Url(foundInstalledApp.Name))
	if storeApp == nil {
		return &installedApp, nil
	}

	installedApp.CurrentVersion = storeApp.CurrentVersion
	return &installedApp, nil
}
