package caddy

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
