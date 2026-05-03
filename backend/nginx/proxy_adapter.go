package nginx

import "github.com/syncloud/platform/config"

type ProxyConfigAdapter struct {
	customProxy *config.CustomProxy
}

func NewProxyConfigAdapter(customProxy *config.CustomProxy) *ProxyConfigAdapter {
	return &ProxyConfigAdapter{customProxy: customProxy}
}

func (a *ProxyConfigAdapter) Proxies() ([]ProxyEntry, error) {
	entries, err := a.customProxy.List()
	if err != nil {
		return nil, err
	}
	result := make([]ProxyEntry, len(entries))
	for i, e := range entries {
		result[i] = ProxyEntry{Name: e.Name, Host: e.Host, Port: e.Port, Https: e.Https, Authelia: e.Authelia}
	}
	return result, nil
}
