package caddy

type ProxyEntry struct {
	Name     string
	Host     string
	Port     int
	Https    bool
	Authelia bool
}
