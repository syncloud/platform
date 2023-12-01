package snap

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
)

type ChangesClient struct {
	logger *zap.Logger
	client ChangesHttpClient
}

type ChangesHttpClient interface {
	Get(url string) ([]byte, error)
}

func NewChangesClient(client ChangesHttpClient, logger *zap.Logger) *ChangesClient {
	return &ChangesClient{
		client: client,
		logger: logger,
	}
}

func (s *ChangesClient) Changes() (*model.InstallerStatus, error) {
	s.logger.Info("snap changes")

	bodyBytes, err := s.client.Get("http://unix/v2/changes?select=in-progress")
	if err != nil {
		return nil, err
	}
	var response model.ServerResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}
	if response.Status != "OK" {
		var errorResponse model.ServerError
		err = json.Unmarshal(response.Result, &errorResponse)
		if err != nil {
			s.logger.Error("cannot unmarshal", zap.Error(err))
			return nil, err
		}

		return nil, fmt.Errorf(errorResponse.Message)
	}

	var changesResponse []model.Change
	err = json.Unmarshal(response.Result, &changesResponse)
	if err != nil {
		s.logger.Error("cannot unmarshal", zap.Error(err))
		return nil, err
	}

	if len(changesResponse) > 0 {
		return &model.InstallerStatus{
			IsRunning: true,
			Progress:  changesResponse[0].InstallerProgress(),
		}, nil
	}
	return &model.InstallerStatus{IsRunning: false}, nil
}
