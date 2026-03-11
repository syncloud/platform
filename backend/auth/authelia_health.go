package auth

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type AutheliaHealth struct {
	socketPath string
	logger     *zap.Logger
}

func NewAutheliaHealth(socketPath string, logger *zap.Logger) *AutheliaHealth {
	return &AutheliaHealth{
		socketPath: socketPath,
		logger:     logger,
	}
}

func (h *AutheliaHealth) WaitForReady() error {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", h.socketPath)
			},
		},
		Timeout: 2 * time.Second,
	}
	for i := 0; i < 30; i++ {
		resp, err := client.Get("http://authelia/api/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				h.logger.Info("authelia is ready")
				return nil
			}
		}
		h.logger.Info("waiting for authelia", zap.Int("attempt", i+1))
		time.Sleep(time.Second)
	}
	return fmt.Errorf("authelia not ready after 30s")
}
