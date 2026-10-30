package caddy

type SystemConfig interface {
	ConfigDir() string
	DataDir() string
}
