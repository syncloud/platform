package caddy

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_OneEntry_AutheliaEnabled(t *testing.T) {
	caddy, _, outputDir := newTestCaddy(t, "example.com", []ProxyEntry{
		{Name: "secret", Host: "192.168.1.20", Port: 9090, Authelia: true},
	})
	assert.NoError(t, caddy.InitCustomProxyConfig())
	assertGolden(t, path.Join(outputDir, "custom-proxy.caddy"), "custom-proxy.one_authelia.caddy")
}
