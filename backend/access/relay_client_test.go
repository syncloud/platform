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
	domain string
	token  *string
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
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{"name.syncloud.it", &token}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())

	err := client.Enable()
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
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{"name.syncloud.it", &token}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())

	assert.Nil(t, client.Enable())
	assert.Nil(t, client.Enable())

	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_EnableWithoutTokenFails(t *testing.T) {
	client := NewRelayClient(&relayControlStub{}, &relaySystemConfigStub{t.TempDir(), "../../config"}, &relayUserConfigStub{"name.syncloud.it", nil}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())
	assert.NotNil(t, client.Enable())
}

func TestRelayClient_DisableRemovesConfigAndRestartsToIdle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "frpc.toml")
	assert.Nil(t, os.WriteFile(path, []byte("x"), 0644))
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{dir, "../../config"}, &relayUserConfigStub{"name.syncloud.it", nil}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())

	err := client.Disable()
	assert.Nil(t, err)
	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr))
	assert.Equal(t, []string{RelayService}, control.restarted)
}

func TestRelayClient_DisableWithoutConfigIsNoop(t *testing.T) {
	control := &relayControlStub{}
	client := NewRelayClient(control, &relaySystemConfigStub{t.TempDir(), "../../config"}, &relayUserConfigStub{"name.syncloud.it", nil}, &relayRedirectStub{"syncloud.it"}, relayRunningClient("name.syncloud.it"), zap.NewNop())
	assert.Nil(t, client.Disable())
	assert.Empty(t, control.restarted)
}
