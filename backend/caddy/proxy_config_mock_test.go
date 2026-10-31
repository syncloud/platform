package caddy

type ProxyConfigMock struct {
	entries []ProxyEntry
}

func (p *ProxyConfigMock) Proxies() ([]ProxyEntry, error) {
	return p.entries, nil
}
