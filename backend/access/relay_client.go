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

	// must match the suffix redirect strips when it attributes tunnel traffic
	// and authorises the proxy's port
	SmtpProxySuffix = "-smtp"
)

type RelayControl interface {
	RestartService(service string) error
}

type RelaySystemConfig interface {
	DataDir() string
	ConfigDir() string
}

type frpcConfig struct {
	Server        string
	Token         string
	AdminSocket   string
	Domain        string
	Web           bool
	LocalPort     int
	Mail          bool
	MailLocalPort int
}

type RelayUserConfig interface {
	GetDeviceDomain() string
	GetDomainUpdateToken() *string
	IsMailRelayEnabled() bool
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

	connectAttempts int
}

func NewRelayClient(control RelayControl, systemConfig RelaySystemConfig, userConfig RelayUserConfig, redirect RelayRedirectConfig, client *http.Client, logger *zap.Logger) *RelayClient {
	return &RelayClient{
		control:      control,
		systemConfig: systemConfig,
		userConfig:   userConfig,
		redirect:     redirect,
		client:       client,
		logger:       logger,

		connectAttempts: relayConnectAttempts,
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

// Apply brings the tunnel to the state the device's settings ask for. The web
// proxy carries app traffic when the relay is on; the smtp proxy carries
// inbound mail when the mail relay is on. Either can be present without the
// other, so a device that only wants mail still gets a tunnel.
func (c *RelayClient) Apply(relayEnabled bool) error {
	mail := c.userConfig.IsMailRelayEnabled()

	if !relayEnabled && !mail {
		return c.Disable()
	}

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
	settings := frpcConfig{
		Server:        server,
		Token:         *token,
		AdminSocket:   c.adminSocketPath(),
		Domain:        domain,
		Web:           relayEnabled,
		LocalPort:     config.WebAccessPort,
		Mail:          mail,
		MailLocalPort: config.MailInboundPort,
	}
	var content bytes.Buffer
	if err := tmpl.Execute(&content, settings); err != nil {
		return err
	}

	expected := c.expectedProxies(domain, relayEnabled, mail)
	if c.currentConfig() == content.String() && c.proxiesRunning(expected) {
		c.logger.Info("relay already connected, skipping restart", zap.Strings("proxies", expected))
		return nil
	}

	c.logger.Info("applying relay",
		zap.String("server", server), zap.Strings("proxies", expected))
	if err := os.WriteFile(c.configPath(), content.Bytes(), 0644); err != nil {
		return err
	}
	if err := c.control.RestartService(RelayService); err != nil {
		return err
	}
	return c.waitConnected(expected)
}

func (c *RelayClient) expectedProxies(domain string, web bool, mail bool) []string {
	var proxies []string
	if web {
		proxies = append(proxies, domain)
	}
	if mail {
		proxies = append(proxies, domain+SmtpProxySuffix)
	}
	return proxies
}

func (c *RelayClient) currentConfig() string {
	content, err := os.ReadFile(c.configPath())
	if err != nil {
		return ""
	}
	return string(content)
}

func (c *RelayClient) waitConnected(proxies []string) error {
	for attempt := 0; attempt < c.connectAttempts; attempt++ {
		if c.proxiesRunning(proxies) {
			c.logger.Info("relay tunnel connected", zap.Strings("proxies", proxies))
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("relay tunnel did not come up for %v", proxies)
}

func (c *RelayClient) proxiesRunning(proxies []string) bool {
	running := c.runningProxies()
	for _, name := range proxies {
		if !running[name] {
			return false
		}
	}
	return true
}

func (c *RelayClient) runningProxies() map[string]bool {
	result := map[string]bool{}
	resp, err := c.client.Get(relayAdminUrl)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return result
	}
	var status map[string][]relayProxyStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return result
	}
	for _, group := range status {
		for _, p := range group {
			if p.Status == "running" {
				result[p.Name] = true
			}
		}
	}
	return result
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
