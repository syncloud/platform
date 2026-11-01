package caddy

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_NoEntries(t *testing.T) {
	caddy, systemd, outputDir := newTestCaddy(t, "example.com", nil)
	assert.NoError(t, caddy.InitCustomProxyConfig())
	assertGolden(t, path.Join(outputDir, "custom-proxy.caddy"), "custom-proxy.empty.caddy")
	assert.Equal(t, "", systemd.reloadedService, "InitCustomProxyConfig should not reload")
}
