package nginx

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SystemdMock struct {
}

func (s *SystemdMock) ReloadService(_ string) error {
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

func TestSubstitution(t *testing.T) {

	outputDir := t.TempDir()

	configDir := path.Join("..", "..", "config")
	systemd := &SystemdMock{}
	systemConfig := &SystemConfigMock{configDir: configDir, dataDir: outputDir}
	userConfig := &UserConfigMock{"example.com"}
	nginx := New(systemd, systemConfig, userConfig)
	err := nginx.InitConfig()
	assert.Nil(t, err)
	resultFile := path.Join(outputDir, "nginx.conf")

	contents, err := os.ReadFile(resultFile)
	assert.Nil(t, err)

	assert.Contains(t, string(contents), "server_name example\\.com;")
	assert.Contains(t, string(contents), "server_name ~^(.*\\.)?(?P<app>.*)\\.example\\.com$;")
}
