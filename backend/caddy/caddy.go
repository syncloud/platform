package caddy

import (
	"bytes"
	"os"
	"path"
	"strings"
	"text/template"
)

const Service = "platform.caddy"

type Caddy struct {
	systemd      Systemd
	systemConfig SystemConfig
	userConfig   UserConfig
	proxyConfig  ProxyConfig
}

func New(systemd Systemd, systemConfig SystemConfig, userConfig UserConfig, proxyConfig ProxyConfig) *Caddy {
	return &Caddy{
		systemd:      systemd,
		userConfig:   userConfig,
		systemConfig: systemConfig,
		proxyConfig:  proxyConfig,
	}
}

func (c *Caddy) ReloadPublic() error {
	return c.systemd.ReloadService(Service)
}

func (c *Caddy) InitConfig() error {
	domain := c.userConfig.GetDeviceDomain()
	configDir := c.systemConfig.ConfigDir()
	templateFile, err := os.ReadFile(path.Join(configDir, "caddy", "Caddyfile"))
	if err != nil {
		return err
	}
	tmpl := string(templateFile)
	tmpl = strings.ReplaceAll(tmpl, "{{ domain_regex }}", strings.ReplaceAll(domain, ".", "\\."))
	tmpl = strings.ReplaceAll(tmpl, "{{ domain }}", domain)
	tmpl = strings.ReplaceAll(tmpl, "{{ auth_host }}", strings.TrimPrefix(c.userConfig.Url("auth"), "https://"))
	return os.WriteFile(path.Join(c.systemConfig.DataDir(), "Caddyfile"), []byte(tmpl), 0644)
}

func (c *Caddy) InitCustomProxyConfig() error {
	return c.writeCustomProxyConfig()
}

func (c *Caddy) ReloadCustomProxy() error {
	err := c.writeCustomProxyConfig()
	if err != nil {
		return err
	}
	return c.systemd.ReloadService(Service)
}

func (c *Caddy) writeCustomProxyConfig() error {
	domain := c.userConfig.GetDeviceDomain()
	configDir := c.systemConfig.ConfigDir()
	templateFile := path.Join(configDir, "caddy", "custom-proxy.caddy")

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	entries, err := c.proxyConfig.Proxies()
	if err != nil {
		return err
	}

	serverEntries := make([]customProxyServerEntry, len(entries))
	for i, e := range entries {
		scheme := "http"
		if e.Https {
			scheme = "https"
		}
		serverEntries[i] = customProxyServerEntry{
			Name:       e.Name,
			ServerName: e.Name + "." + domain,
			Host:       e.Host,
			Port:       e.Port,
			Scheme:     scheme,
			Authelia:   e.Authelia,
		}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, customProxyTemplateData{Entries: serverEntries})
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(c.systemConfig.DataDir(), "custom-proxy.caddy"), buf.Bytes(), 0644)
}
