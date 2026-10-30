package caddy

type customProxyServerEntry struct {
	Name       string
	ServerName string
	Host       string
	Port       int
	Scheme     string
	Authelia   bool
}
