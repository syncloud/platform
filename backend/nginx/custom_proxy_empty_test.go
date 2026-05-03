package nginx

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_NoEntries(t *testing.T) {
	nginx, systemd, outputDir := newTestNginx(t, "example.com", nil)
	assert.NoError(t, nginx.InitCustomProxyConfig())
	assertGolden(t, path.Join(outputDir, "custom-proxy.conf"), "custom-proxy.empty.conf")
	assert.Equal(t, "", systemd.reloadedService, "InitCustomProxyConfig should not reload")
}
