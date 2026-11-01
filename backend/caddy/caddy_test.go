package caddy

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestCaddy(t *testing.T, domain string, entries []ProxyEntry) (*Caddy, *SystemdMock, string) {
	t.Helper()
	outputDir := t.TempDir()
	configDir := path.Join("..", "..", "config")
	systemd := &SystemdMock{}
	systemConfig := &SystemConfigMock{configDir: configDir, dataDir: outputDir}
	userConfig := &UserConfigMock{deviceDomain: domain}
	proxyConfig := &ProxyConfigMock{entries: entries}
	return New(systemd, systemConfig, userConfig, proxyConfig), systemd, outputDir
}

func assertGolden(t *testing.T, generatedPath, goldenName string) {
	t.Helper()
	actual, err := os.ReadFile(generatedPath)
	assert.NoError(t, err)
	goldenPath := path.Join("testdata", goldenName)
	if os.Getenv("UPDATE_GOLDENS") == "1" {
		assert.NoError(t, os.MkdirAll("testdata", 0755))
		assert.NoError(t, os.WriteFile(goldenPath, actual, 0644))
		return
	}
	expected, err := os.ReadFile(goldenPath)
	assert.NoError(t, err, "golden missing: run with UPDATE_GOLDENS=1 to create %s", goldenPath)
	assert.Equal(t, string(expected), string(actual), "output differs from %s", goldenPath)
}
