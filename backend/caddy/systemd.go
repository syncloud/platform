package caddy

type Systemd interface {
	ReloadService(service string) error
}
