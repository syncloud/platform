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
	domain     string
	token      *string
	mailSocket *string
}

func (c *relayUserConfigStub) GetMailInboundSocket() *string {
	return c.mailSocket
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

func TestRelayClient_EnableWritesConfigAndRestartsAndVerifies(t *testing.T) {
	dir := t.TempDir()
	token := "tok123"
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it", token: &token}, &relayRedirectStub{"syncloud.it"}, relayRunningProxies("name.syncloud.it", "name.syncloud.it-smtp"), zap.NewNop())

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
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it", token: &token}, &relayRedirectStub{"syncloud.it"}, relayRunningProxies("name.syncloud.it", "name.syncloud.it-smtp"), zap.NewNop())

	assert.Nil(t, client.Apply(true))
	assert.Nil(t, client.Apply(true))

	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_EnableWithoutTokenFails(t *testing.T) {
	client := NewRelayClient(&relayControlStub{}, &relaySystemConfigStub{t.TempDir(), "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it"}, &relayRedirectStub{"syncloud.it"}, relayRunningProxies("name.syncloud.it", "name.syncloud.it-smtp"), zap.NewNop())
	assert.NotNil(t, client.Apply(true))
}

func TestRelayClient_DisableRemovesConfigAndRestartsToIdle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "frpc.toml")
	assert.Nil(t, os.WriteFile(path, []byte("x"), 0644))
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it"}, &relayRedirectStub{"syncloud.it"}, relayRunningProxies("name.syncloud.it", "name.syncloud.it-smtp"), zap.NewNop())

	err := client.Disable()
	assert.Nil(t, err)
	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr))
	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_DisableWithoutConfigIsNoop(t *testing.T) {
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{t.TempDir(), "../../config"}, &relayUserConfigStub{domain: "name.syncloud.it"}, &relayRedirectStub{"syncloud.it"}, relayRunningProxies("name.syncloud.it", "name.syncloud.it-smtp"), zap.NewNop())
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

func mailClient(t *testing.T, dir string, socket *string, running ...string) (*RelayClient, *relayControlStub) {
	t.Helper()
	token := "tok123"
	control := &relayControlStub{}
	user := &relayUserConfigStub{domain: "name.syncloud.it", token: &token, mailSocket: socket}
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, user,
		&relayRedirectStub{"syncloud.it"}, relayRunningProxies(running...), zap.NewNop())
	return client, control
}

func TestRelayClient_RegisteredMailOpensTheSmtpProxy(t *testing.T) {
	dir := t.TempDir()
	socket := "/var/snap/mail/common/mail.socket"
	client, _ := mailClient(t, dir, &socket, "name.syncloud.it", "name.syncloud.it-smtp")

	assert.Nil(t, client.Apply(true))

	s, err := readConfig(dir)
	assert.Nil(t, err)
	assert.Contains(t, s, `type = "https"`)
	assert.Contains(t, s, `name = "name.syncloud.it-smtp"`)
	assert.Contains(t, s, `type = "tcpmux"`)
	assert.Contains(t, s, `multiplexer = "httpconnect"`)
	assert.Contains(t, s, `type = "unix_domain_socket"`)
	assert.Contains(t, s, `unixPath = "/var/snap/mail/common/mail.socket"`)
}

func TestRelayClient_WithoutRegisteredMailThereIsNoSmtpProxy(t *testing.T) {
	dir := t.TempDir()
	client, _ := mailClient(t, dir, nil, "name.syncloud.it")

	assert.Nil(t, client.Apply(true))

	s, err := readConfig(dir)
	assert.Nil(t, err)
	assert.Contains(t, s, `type = "https"`)
	assert.NotContains(t, s, "-smtp")
	assert.NotContains(t, s, `type = "tcpmux"`)
}

func TestRelayClient_RelayOffRemovesTheTunnel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "frpc.toml")
	assert.Nil(t, os.WriteFile(path, []byte("x"), 0644))
	client, _ := mailClient(t, dir, nil)

	assert.Nil(t, client.Apply(false))

	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr))
}

func TestRelayClient_WaitsForBothProxiesBeforeSucceeding(t *testing.T) {
	dir := t.TempDir()
	socket := "/var/snap/mail/common/mail.socket"
	client, _ := mailClient(t, dir, &socket, "name.syncloud.it")
	client.connectAttempts = 2

	assert.NotNil(t, client.Apply(true))
}

func readConfig(dir string) (string, error) {
	content, err := os.ReadFile(filepath.Join(dir, "frpc.toml"))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
