package caddy

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_MixedEntries(t *testing.T) {
	caddy, _, outputDir := newTestCaddy(t, "mydevice.syncloud.it", []ProxyEntry{
		{Name: "open", Host: "192.168.1.10", Port: 8080, Authelia: false},
		{Name: "camera", Host: "10.0.0.100", Port: 8443, Https: true},
		{Name: "secret", Host: "192.168.1.20", Port: 9090, Authelia: true},
	})
	assert.NoError(t, caddy.InitCustomProxyConfig())
	assertGolden(t, path.Join(outputDir, "custom-proxy.caddy"), "custom-proxy.mixed.caddy")
}
