package caddy

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const gatewayBind = "bind unix//var/snap/platform/current/authelia.socket|0777"

func TestPublicConfig_Substitution(t *testing.T) {
	caddy, _, outputDir := newTestCaddy(t, "example.com", nil)
	assert.NoError(t, caddy.InitConfig())
	assertGolden(t, path.Join(outputDir, "Caddyfile"), "Caddyfile.example.com")
}

func TestPublicConfig_ReloadsCaddy(t *testing.T) {
	caddy, systemd, _ := newTestCaddy(t, "example.com", nil)
	assert.NoError(t, caddy.ReloadPublic())
	assert.Equal(t, Service, systemd.reloadedService)
}

func TestPublicConfig_IndexRevalidates(t *testing.T) {
	config := generatePublicConfig(t)
	assert.Equal(t, 2, strings.Count(config, "Cache-Control no-cache"))
}

func TestPublicConfig_AssetsImmutableWithoutIndexFallback(t *testing.T) {
	config := generatePublicConfig(t)
	assets := block(config, "handle /assets/* {")
	assert.Contains(t, assets, `header Cache-Control "public, max-age=31536000, immutable"`)
	assert.NotContains(t, assets, "try_files")
	assert.Equal(t, 2, strings.Count(config, "import spa "))
}

func TestPublicConfig_SecurityHeadersOnEverySite(t *testing.T) {
	config := generatePublicConfig(t)
	assert.Equal(t, 3, strings.Count(config, `Strict-Transport-Security "max-age=31536000; includeSubdomains"`))
	assert.Equal(t, 2, strings.Count(config, "Access-Control-Allow-Origin *"))
}

func TestPublicConfig_NoHttp3(t *testing.T) {
	config := generatePublicConfig(t)
	assert.Equal(t, 2, strings.Count(config, "protocols h1 h2"))
}

func TestPublicConfig_StreamAppsGetDomain(t *testing.T) {
	config := generatePublicConfig(t)
	assert.Contains(t, config, "@stream syncloud_stream example.com")
	assert.Contains(t, config, "upstream unix//var/snap/{l4.syncloud.app}/common/stream.socket")
}

func TestPublicConfig_AuthGatewaySendsThePublicAuthHost(t *testing.T) {
	gateway := strings.SplitN(generatePublicConfig(t), gatewayBind, 2)[1]
	assert.Contains(t, gateway, "header_up Host auth.example.com\n")
	assert.Contains(t, gateway, "header_up X-Forwarded-Host auth.example.com\n")
	assert.Contains(t, gateway, "header_up X-Forwarded-Proto https\n")
}

func TestPublicConfig_AuthGatewayKeepsTheExternalPort(t *testing.T) {
	caddy, _, outputDir := newTestCaddy(t, "example.com", nil)
	port := 10443
	caddy.userConfig.(*UserConfigMock).port = &port
	assert.NoError(t, caddy.InitConfig())
	content, err := os.ReadFile(path.Join(outputDir, "Caddyfile"))
	assert.NoError(t, err)
	gateway := strings.SplitN(string(content), gatewayBind, 2)[1]
	assert.Contains(t, gateway, "header_up Host auth.example.com:10443\n")
	assert.Contains(t, gateway, "header_up X-Forwarded-Host auth.example.com:10443\n")
}

func TestPublicConfig_AuthGatewayLeavesForwardAuthHeadersAlone(t *testing.T) {
	gateway := strings.SplitN(generatePublicConfig(t), gatewayBind, 2)[1]
	assert.NotContains(t, block(gateway, "handle /api/authz/* {"), "header_up")
}

func generatePublicConfig(t *testing.T) string {
	t.Helper()
	caddy, _, outputDir := newTestCaddy(t, "example.com", nil)
	assert.NoError(t, caddy.InitConfig())
	content, err := os.ReadFile(path.Join(outputDir, "Caddyfile"))
	assert.NoError(t, err)
	return string(content)
}

func block(config, opening string) string {
	return strings.SplitN(strings.SplitN(config, opening, 2)[1], "\n\t}", 2)[0]
}
