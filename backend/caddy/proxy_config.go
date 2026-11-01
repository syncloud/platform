package caddy

type ProxyConfig interface {
	Proxies() ([]ProxyEntry, error)
}
