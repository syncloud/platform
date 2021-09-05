package nginx

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type SystemdMock struct {
}

func (s *SystemdMock) ReloadService(service string) error {
	return nil
}

type SystemConfigMock struct {
	configDir string
	dataDir   string
}

func (s *SystemConfigMock) ConfigDir() (*string, error) {
	return &s.configDir, nil
}

func (s *SystemConfigMock) DataDir() (*string, error) {
	return &s.dataDir, nil
}

type UserConfigMock struct {
	deviceDomain string
}

func (u *UserConfigMock) GetDeviceDomain() string {
	return u.deviceDomain
}

func TestSubstitution(t *testing.T) {

	outputDir, err := ioutil.TempDir("", "")
	defer func() { _ = os.Remove(outputDir) }()

	configDir := path.Join("..", "..", "config")
	systemd := &SystemdMock{}
	systemConfig := &SystemConfigMock{configDir: configDir, dataDir: outputDir}
	userConfig := &UserConfigMock{"example.com"}
	nginx := New(systemd, systemConfig, userConfig)
	assert.Nil(t, err)
	err = nginx.InitConfig()
	assert.Nil(t, err)
	resultFile := path.Join(outputDir, "nginx.conf")

	contents, err := ioutil.ReadFile(resultFile)
	assert.Nil(t, err)

	assert.Contains(t, string(contents), "server_name example\\.com;")
	assert.Contains(t, string(contents), "server_name ~^(?P<app>.*)\\.example\\.com$;")
}
