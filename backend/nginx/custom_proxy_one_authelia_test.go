package nginx

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_OneEntry_AutheliaEnabled(t *testing.T) {
	nginx, _, outputDir := newTestNginx(t, "example.com", []ProxyEntry{
		{Name: "secret", Host: "192.168.1.20", Port: 9090, Authelia: true},
	})
	assert.NoError(t, nginx.InitCustomProxyConfig())
	assertGolden(t, path.Join(outputDir, "custom-proxy.conf"), "custom-proxy.one_authelia.conf")
}
