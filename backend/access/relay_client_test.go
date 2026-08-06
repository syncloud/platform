package access

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type relayControlStub struct {
	restarted []string
}

func (c *relayControlStub) RestartService(s string) error {
	c.restarted = append(c.restarted, s)
	return nil
}

type relaySystemConfigStub struct {
	dataDir   string
	configDir string
}

func (c *relaySystemConfigStub) DataDir() string {
	return c.dataDir
}

func (c *relaySystemConfigStub) ConfigDir() string {
	return c.configDir
}

type relayUserConfigStub struct {
	mailRelay bool
	smtpPort  *int
	domain    string
	token     *string
}

func (s *relayUserConfigStub) IsMailRelayEnabled() bool {
	return s.mailRelay
}

func (s *relayUserConfigStub) GetMailSmtpPort() *int {
	return s.smtpPort
}

func (c *relayUserConfigStub) GetDeviceDomain() string {
	return c.domain
}

func (c *relayUserConfigStub) GetDomainUpdateToken() *string {
	return c.token
}

type relayRedirectStub struct {
	domain string
}

func (c *relayRedirectStub) Domain() string {
	return c.domain
}

type relayStatusTransport struct {
	body string
}

func (t relayStatusTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(t.body)), Header: make(http.Header)}, nil
}

func relayRunningClient(domain string) *http.Client {
	body := fmt.Sprintf(`{"https":[{"name":"%s","status":"running"}]}`, domain)
	return &http.Client{Transport: relayStatusTransport{body}}
}

func TestRelayClient_EnableWritesConfigAndRestartsAndVerifies(t *testing.T) {
	dir := t.TempDir()
	token := "tok123"
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it", token: &token}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())

	err := client.Apply(true)
	assert.Nil(t, err)

	content, err := os.ReadFile(filepath.Join(dir, "frpc.toml"))
	assert.Nil(t, err)
	s := string(content)
	assert.Contains(t, s, `serverAddr = "relay.syncloud.it"`)
	assert.Contains(t, s, `metadatas.token = "tok123"`)
	assert.Contains(t, s, `name = "name.syncloud.it"`)
	assert.Contains(t, s, `customDomains = ["name.syncloud.it", "*.name.syncloud.it"]`)
	assert.Contains(t, s, fmt.Sprintf(`webServer.unixSocket = "%s"`, filepath.Join(dir, "frpc-admin.sock")))
	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_EnableIdempotentSkipsRestartWhenConnected(t *testing.T) {
	dir := t.TempDir()
	token := "tok123"
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it", token: &token}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())

	assert.Nil(t, client.Apply(true))
	assert.Nil(t, client.Apply(true))

	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_EnableWithoutTokenFails(t *testing.T) {
	client := NewRelayClient(&relayControlStub{}, &relaySystemConfigStub{t.TempDir(), "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it"}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())
	assert.NotNil(t, client.Apply(true))
}

func TestRelayClient_DisableRemovesConfigAndRestartsToIdle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "frpc.toml")
	assert.Nil(t, os.WriteFile(path, []byte("x"), 0644))
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it"}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())

	err := client.Disable()
	assert.Nil(t, err)
	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr))
	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_DisableWithoutConfigIsNoop(t *testing.T) {
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{t.TempDir(), "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it"}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())
	assert.Nil(t, client.Disable())
	assert.Empty(t, control.restarted)
}

func relayRunningProxies(names ...string) *http.Client {
	var parts []string
	for _, n := range names {
		parts = append(parts, fmt.Sprintf(`{"name":"%s","status":"running"}`, n))
	}
	body := fmt.Sprintf(`{"all":[%s]}`, strings.Join(parts, ","))
	return &http.Client{Transport: relayStatusTransport{body}}
}

func mailClient(t *testing.T, dir string, relay bool, mail bool, port *int, running ...string) (*RelayClient, *relayControlStub) {
	t.Helper()
	token := "tok123"
	control := &relayControlStub{}
	user := &relayUserConfigStub{domain: "name.syncloud.it", token: &token, mailRelay: mail, smtpPort: port}
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, user,
		&relayRedirectStub{"syncloud.it"}, relayRunningProxies(running...), zap.NewNop())
	return client, control
}

func TestRelayClient_MailOnlyTunnelHasNoWebProxy(t *testing.T) {
	dir := t.TempDir()
	port := 20005
	client, control := mailClient(t, dir, false, true, &port, "name.syncloud.it-smtp")

	assert.Nil(t, client.Apply(false))

	s := readConfig(t, dir)
	assert.NotContains(t, s, `type = "https"`)
	assert.Contains(t, s, `name = "name.syncloud.it-smtp"`)
	assert.Contains(t, s, `type = "tcp"`)
	assert.Contains(t, s, "remotePort = 20005")
	assert.Contains(t, s, "localPort = 10025")
	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_BothProxiesWhenRelayAndMailAreOn(t *testing.T) {
	dir := t.TempDir()
	port := 20005
	client, _ := mailClient(t, dir, true, true, &port, "name.syncloud.it", "name.syncloud.it-smtp")

	assert.Nil(t, client.Apply(true))

	s := readConfig(t, dir)
	assert.Contains(t, s, `type = "https"`)
	assert.Contains(t, s, `type = "tcp"`)
	assert.Contains(t, s, "remotePort = 20005")
}

func TestRelayClient_NoMailProxyWithoutAnAllocatedPort(t *testing.T) {
	dir := t.TempDir()
	client, _ := mailClient(t, dir, true, true, nil, "name.syncloud.it")

	assert.Nil(t, client.Apply(true))

	s := readConfig(t, dir)
	assert.Contains(t, s, `type = "https"`)
	assert.NotContains(t, s, `type = "tcp"`)
}

func TestRelayClient_MailRelayOnButNoPortAndNoRelayIsIdle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "frpc.toml")
	assert.Nil(t, os.WriteFile(path, []byte("x"), 0644))
	client, control := mailClient(t, dir, false, true, nil)

	assert.Nil(t, client.Apply(false))

	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr))
	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_BothOffRemovesTheTunnel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "frpc.toml")
	assert.Nil(t, os.WriteFile(path, []byte("x"), 0644))
	client, _ := mailClient(t, dir, false, false, nil)

	assert.Nil(t, client.Apply(false))

	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr))
}

func TestRelayClient_WaitsForBothProxiesBeforeSucceeding(t *testing.T) {
	dir := t.TempDir()
	port := 20005
	// only the web proxy comes up, the mail one never does
	client, _ := mailClient(t, dir, true, true, &port, "name.syncloud.it")
	client.connectAttempts = 2

	assert.NotNil(t, client.Apply(true))
}

func readConfig(t *testing.T, dir string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(dir, "frpc.toml"))
	assert.Nil(t, err)
	return string(content)
}
