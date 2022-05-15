package auth

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"os/exec"
)

type SystemPasswordChanger struct {
	logger *zap.Logger
}

type PasswordChanger interface {
	Change(password string) error
}

func NewSystemPassword(logger *zap.Logger) *SystemPasswordChanger {
	return &SystemPasswordChanger{logger: logger}
}

func (s *SystemPasswordChanger) Change(password string) error {
	cmd := exec.Command("chpasswd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer func() { _ = stdin.Close() }()
		io.WriteString(stdin, fmt.Sprintf("root:%s\n", password))
	}()
	out, err := cmd.CombinedOutput()
	s.logger.Info("chpasswd", zap.ByteString("output", out))
	return err
}
