package nginx

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_OneEntry_NoAuthelia(t *testing.T) {
	nginx, _, outputDir := newTestNginx(t, "example.com", []ProxyEntry{
		{Name: "myapp", Host: "192.168.1.10", Port: 8080},
	})
	assert.NoError(t, nginx.InitCustomProxyConfig())
	assertGolden(t, path.Join(outputDir, "custom-proxy.conf"), "custom-proxy.one_no_authelia.conf")
}
