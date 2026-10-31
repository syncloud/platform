package caddy

type UserConfig interface {
	GetDeviceDomain() string
	Url(app string) string
}
