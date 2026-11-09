package smart

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"strings"
)

const SMARTCTL = "/snap/platform/current/smartmontools/bin/smartctl.sh"

type Smartctl struct {
	executor cli.Executor
	logger   *zap.Logger
}

func NewSmartctl(executor cli.Executor, logger *zap.Logger) *Smartctl {
	return &Smartctl{
		executor: executor,
		logger:   logger,
	}
}

func (s *Smartctl) Health(device string) model.SmartHealth {
	output, _, err := s.detect(device)
	if err != nil {
		s.logger.Warn("smartctl failed", zap.String("device", device), zap.Error(err))
		return model.SmartHealth{Message: err.Error()}
	}
	if !output.Supported() {
		return model.SmartHealth{Message: output.Message()}
	}
	return output.Health()
}

func (s *Smartctl) StartShortTest(device string) error {
	output, deviceType, err := s.detect(device)
	if err != nil {
		return err
	}
	if !output.Supported() {
		return errors.New(output.Message())
	}
	if !output.SelfTestSupported() {
		return errors.New(SelfTestNotSupported)
	}
	result, err := s.run(append(deviceType, "-t", "short", device)...)
	if err != nil {
		return err
	}
	if !result.Succeeded() {
		return errors.New(result.Failure())
	}
	return nil
}

func (s *Smartctl) detect(device string) (*Output, []string, error) {
	output, err := s.run("-a", device)
	if err != nil {
		return nil, nil, err
	}
	if output.Supported() {
		return output, nil, nil
	}
	deviceType := []string{"-d", "sat"}
	sat, err := s.run(append(deviceType, "-a", device)...)
	if err == nil && sat.Supported() {
		return sat, deviceType, nil
	}
	return output, nil, nil
}

func (s *Smartctl) run(args ...string) (*Output, error) {
	raw, execErr := s.executor.CombinedOutput(SMARTCTL, append([]string{"-j"}, args...)...)
	var output Output
	if err := json.Unmarshal(raw, &output); err != nil {
		if execErr != nil {
			return nil, execErr
		}
		return nil, fmt.Errorf("unable to parse smartctl output: %s", strings.TrimSpace(string(raw)))
	}
	return &output, nil
}
