package nginx

import "github.com/syncloud/platform/config"

type ProxyConfigAdapter struct {
	userConfig *config.UserConfig
}

func NewProxyConfigAdapter(userConfig *config.UserConfig) *ProxyConfigAdapter {
	return &ProxyConfigAdapter{userConfig: userConfig}
}

func (a *ProxyConfigAdapter) Proxies() ([]ProxyEntry, error) {
	entries, err := a.userConfig.CustomProxies()
	if err != nil {
		return nil, err
	}
	result := make([]ProxyEntry, len(entries))
	for i, e := range entries {
		result[i] = ProxyEntry{Name: e.Name, Host: e.Host, Port: e.Port, Https: e.Https, Authelia: e.Authelia}
	}
	return result, nil
}
