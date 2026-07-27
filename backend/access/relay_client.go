package access

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/syncloud/platform/config"
	"go.uber.org/zap"
)

const (
	RelayService         = "platform.frpc"
	relayAdminSocket     = "frpc-admin.sock"
	relayAdminUrl        = "http://unix/api/status"
	relayConnectAttempts = 30
)

type RelayControl interface {
	RestartService(service string) error
}

type RelaySystemConfig interface {
	DataDir() string
	ConfigDir() string
}

type frpcConfig struct {
	Server      string
	Token       string
	AdminSocket string
	Domain      string
	LocalPort   int
}

type RelayUserConfig interface {
	GetDeviceDomain() string
	GetDomainUpdateToken() *string
}

type RelayRedirectConfig interface {
	Domain() string
}

type RelayClient struct {
	control      RelayControl
	systemConfig RelaySystemConfig
	userConfig   RelayUserConfig
	redirect     RelayRedirectConfig
	client       *http.Client
	logger       *zap.Logger
}

func NewRelayClient(control RelayControl, systemConfig RelaySystemConfig, userConfig RelayUserConfig, redirect RelayRedirectConfig, client *http.Client, logger *zap.Logger) *RelayClient {
	return &RelayClient{
		control:      control,
		systemConfig: systemConfig,
		userConfig:   userConfig,
		redirect:     redirect,
		client:       client,
		logger:       logger,
	}
}

func NewFrpcAdminClient(dataDir string) *http.Client {
	socket := filepath.Join(dataDir, relayAdminSocket)
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", socket)
			},
		},
	}
}

func (c *RelayClient) configPath() string {
	return filepath.Join(c.systemConfig.DataDir(), "frpc.toml")
}

func (c *RelayClient) adminSocketPath() string {
	return filepath.Join(c.systemConfig.DataDir(), relayAdminSocket)
}

func (c *RelayClient) Enable() error {
	domain := c.userConfig.GetDeviceDomain()
	token := c.userConfig.GetDomainUpdateToken()
	if token == nil {
		return fmt.Errorf("domain update token is not available")
	}
	server := fmt.Sprintf("relay.%s", c.redirect.Domain())
	tmpl, err := template.ParseFiles(filepath.Join(c.systemConfig.ConfigDir(), "frp", "frpc.toml"))
	if err != nil {
		return err
	}
	var content bytes.Buffer
	err = tmpl.Execute(&content, frpcConfig{
		Server:      server,
		Token:       *token,
		AdminSocket: c.adminSocketPath(),
		Domain:      domain,
		LocalPort:   config.WebAccessPort,
	})
	if err != nil {
		return err
	}

	if c.currentConfig() == content.String() && c.proxyRunning(domain) {
		c.logger.Info("relay already connected, skipping restart", zap.String("domain", domain))
		return nil
	}

	c.logger.Info("enabling relay", zap.String("server", server), zap.String("domain", domain))
	if err := os.WriteFile(c.configPath(), content.Bytes(), 0644); err != nil {
		return err
	}
	if err := c.control.RestartService(RelayService); err != nil {
		return err
	}
	return c.waitConnected(domain)
}

func (c *RelayClient) currentConfig() string {
	content, err := os.ReadFile(c.configPath())
	if err != nil {
		return ""
	}
	return string(content)
}

func (c *RelayClient) waitConnected(domain string) error {
	for attempt := 0; attempt < relayConnectAttempts; attempt++ {
		if c.proxyRunning(domain) {
			c.logger.Info("relay tunnel connected", zap.String("domain", domain))
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("relay tunnel did not come up for %s", domain)
}

func (c *RelayClient) proxyRunning(domain string) bool {
	resp, err := c.client.Get(relayAdminUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	var status map[string][]relayProxyStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return false
	}
	for _, proxies := range status {
		for _, p := range proxies {
			if p.Name == domain && p.Status == "running" {
				return true
			}
		}
	}
	return false
}

func (c *RelayClient) Disable() error {
	if _, err := os.Stat(c.configPath()); os.IsNotExist(err) {
		return nil
	}
	c.logger.Info("disabling relay")
	if err := os.Remove(c.configPath()); err != nil {
		return err
	}
	return c.control.RestartService(RelayService)
}
