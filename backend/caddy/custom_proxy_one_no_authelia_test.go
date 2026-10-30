package caddy

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_OneEntry_NoAuthelia(t *testing.T) {
	caddy, _, outputDir := newTestCaddy(t, "example.com", []ProxyEntry{
		{Name: "myapp", Host: "192.168.1.10", Port: 8080},
	})
	assert.NoError(t, caddy.InitCustomProxyConfig())
	assertGolden(t, path.Join(outputDir, "custom-proxy.caddy"), "custom-proxy.one_no_authelia.caddy")
}
