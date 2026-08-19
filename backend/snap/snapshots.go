package snap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
)

const RetentionDisabled = "no"

type SnapshotsClient interface {
	Get(url string) ([]byte, error)
	Post(url, bodyType string, body io.Reader) (*http.Response, error)
	Put(url, bodyType string, body io.Reader) (*http.Response, error)
}

type Snapshots struct {
	client SnapshotsClient
	logger *zap.Logger
}

func NewSnapshots(client SnapshotsClient, logger *zap.Logger) *Snapshots {
	return &Snapshots{
		client: client,
		logger: logger,
	}
}

func (s *Snapshots) DisableAutomatic() error {
	retention, err := s.retention()
	if err != nil {
		return err
	}
	if retention == RetentionDisabled {
		return nil
	}
	s.logger.Info("disabling automatic snapshots", zap.String("was", retention))
	requestJson, err := json.Marshal(model.NewSnapshotsConf(RetentionDisabled))
	if err != nil {
		return err
	}
	resp, err := s.client.Put("http://unix/v2/snaps/system/conf", "application/json", bytes.NewBuffer(requestJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("cannot disable automatic snapshots: %s", string(body))
	}
	return nil
}

func (s *Snapshots) ForgetAuto() error {
	sets, err := s.List()
	if err != nil {
		return err
	}
	for _, set := range sets {
		if !set.IsAuto() {
			continue
		}
		err = s.Forget(set.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Snapshots) List() ([]model.SnapshotSet, error) {
	body, err := s.client.Get("http://unix/v2/snapshots")
	if err != nil {
		return nil, err
	}
	var response model.SnapshotSetsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}
	return response.Result, nil
}

func (s *Snapshots) Forget(set int) error {
	s.logger.Info("forgetting snapshot", zap.Int("set", set))
	requestJson, err := json.Marshal(model.SnapshotsRequest{Action: "forget", Set: set})
	if err != nil {
		return err
	}
	resp, err := s.client.Post("http://unix/v2/snapshots", "application/json", bytes.NewBuffer(requestJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("cannot forget snapshot %d: %s", set, string(body))
	}
	return nil
}

func (s *Snapshots) retention() (string, error) {
	body, err := s.client.Get(fmt.Sprintf("http://unix/v2/snaps/system/conf?keys=%s", model.RetentionKey))
	if err != nil {
		return "", err
	}
	var response model.ConfResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return "", err
	}
	return response.Value(model.RetentionKey), nil
}
