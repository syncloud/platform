package caddy

type SystemdMock struct {
	reloadedService string
}

func (s *SystemdMock) ReloadService(service string) error {
	s.reloadedService = service
	return nil
}
