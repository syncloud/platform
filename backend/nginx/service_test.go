package nginx

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SystemdMock struct {
	reloadedService string
}

func (s *SystemdMock) ReloadService(service string) error {
	s.reloadedService = service
	return nil
}

type SystemConfigMock struct {
	configDir string
	dataDir   string
}

func (s *SystemConfigMock) ConfigDir() string {
	return s.configDir
}

func (s *SystemConfigMock) DataDir() string {
	return s.dataDir
}

type UserConfigMock struct {
	deviceDomain string
}

func (u *UserConfigMock) GetDeviceDomain() string {
	return u.deviceDomain
}

type ProxyConfigMock struct {
	entries []ProxyEntry
}

func (p *ProxyConfigMock) Proxies() ([]ProxyEntry, error) {
	return p.entries, nil
}

func TestSubstitution(t *testing.T) {
	outputDir := t.TempDir()
	configDir := path.Join("..", "..", "config")
	systemd := &SystemdMock{}
	systemConfig := &SystemConfigMock{configDir: configDir, dataDir: outputDir}
	userConfig := &UserConfigMock{"example.com"}
	proxyConfig := &ProxyConfigMock{}
	nginx := New(systemd, systemConfig, userConfig, proxyConfig)
	err := nginx.InitConfig()
	assert.Nil(t, err)

	contents, err := os.ReadFile(path.Join(outputDir, "nginx.conf"))
	assert.Nil(t, err)

	assert.Contains(t, string(contents), "server_name example.com;")
	assert.Contains(t, string(contents), "server_name ~^(.*\\.)?(?P<app>.*)\\.example\\.com$;")
	assert.Contains(t, string(contents), "@custom_proxy")
}

func TestCustomProxy_ZeroEntries(t *testing.T) {
	outputDir := t.TempDir()
	configDir := path.Join("..", "..", "config")
	systemd := &SystemdMock{}
	systemConfig := &SystemConfigMock{configDir: configDir, dataDir: outputDir}
	userConfig := &UserConfigMock{"example.com"}
	proxyConfig := &ProxyConfigMock{}
	nginx := New(systemd, systemConfig, userConfig, proxyConfig)

	err := nginx.InitCustomProxyConfig()
	assert.Nil(t, err)

	contents, err := os.ReadFile(path.Join(outputDir, "custom-proxy.conf"))
	assert.Nil(t, err)
	text := string(contents)

	assert.Contains(t, text, "return 502")
	assert.NotContains(t, text, "proxy_pass http://")
	assert.Equal(t, 1, strings.Count(text, "listen unix:"), "should have only the default server block")
	assert.Equal(t, "platform.nginx-custom-proxy", systemd.reloadedService)
}

func TestCustomProxy_OneEntry(t *testing.T) {
	outputDir := t.TempDir()
	configDir := path.Join("..", "..", "config")
	systemd := &SystemdMock{}
	systemConfig := &SystemConfigMock{configDir: configDir, dataDir: outputDir}
	userConfig := &UserConfigMock{"example.com"}
	proxyConfig := &ProxyConfigMock{entries: []ProxyEntry{
		{Name: "myapp", Host: "192.168.1.10", Port: 8080},
	}}
	nginx := New(systemd, systemConfig, userConfig, proxyConfig)

	err := nginx.InitCustomProxyConfig()
	assert.Nil(t, err)

	contents, err := os.ReadFile(path.Join(outputDir, "custom-proxy.conf"))
	assert.Nil(t, err)
	text := string(contents)

	assert.Contains(t, text, "server_name myapp.example.com;")
	assert.Contains(t, text, "proxy_pass http://192.168.1.10:8080;")
	assert.Contains(t, text, "X-Syncloud-Custom-Proxy")
	assert.Equal(t, 2, strings.Count(text, "listen unix:"), "should have default + 1 custom server block")
	assert.Equal(t, "platform.nginx-custom-proxy", systemd.reloadedService)
}

func TestCustomProxy_TwoEntries(t *testing.T) {
	outputDir := t.TempDir()
	configDir := path.Join("..", "..", "config")
	systemd := &SystemdMock{}
	systemConfig := &SystemConfigMock{configDir: configDir, dataDir: outputDir}
	userConfig := &UserConfigMock{"mydevice.syncloud.it"}
	proxyConfig := &ProxyConfigMock{entries: []ProxyEntry{
		{Name: "nas", Host: "192.168.1.50", Port: 5000},
		{Name: "camera", Host: "10.0.0.100", Port: 8443},
	}}
	nginx := New(systemd, systemConfig, userConfig, proxyConfig)

	err := nginx.InitCustomProxyConfig()
	assert.Nil(t, err)

	contents, err := os.ReadFile(path.Join(outputDir, "custom-proxy.conf"))
	assert.Nil(t, err)
	text := string(contents)

	assert.Contains(t, text, "server_name nas.mydevice.syncloud.it;")
	assert.Contains(t, text, "proxy_pass http://192.168.1.50:5000;")
	assert.Contains(t, text, "server_name camera.mydevice.syncloud.it;")
	assert.Contains(t, text, "proxy_pass http://10.0.0.100:8443;")
	assert.Equal(t, 3, strings.Count(text, "listen unix:"), "should have default + 2 custom server blocks")
	assert.Equal(t, "platform.nginx-custom-proxy", systemd.reloadedService)
}
