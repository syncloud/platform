package nginx

import (
	"bytes"
	"os"
	"path"
	"strings"
	"text/template"
)

type Systemd interface {
	ReloadService(service string) error
}

type SystemConfig interface {
	ConfigDir() string
	DataDir() string
}

type UserConfig interface {
	GetDeviceDomain() string
}

type ProxyEntry struct {
	Name string
	Host string
	Port int
}

type ProxyConfig interface {
	Proxies() ([]ProxyEntry, error)
}

type Nginx struct {
	systemd      Systemd
	systemConfig SystemConfig
	userConfig   UserConfig
	proxyConfig  ProxyConfig
}

func New(systemd Systemd, systemConfig SystemConfig, userConfig UserConfig, proxyConfig ProxyConfig) *Nginx {
	return &Nginx{
		systemd:      systemd,
		userConfig:   userConfig,
		systemConfig: systemConfig,
		proxyConfig:  proxyConfig,
	}
}

func (n *Nginx) ReloadPublic() error {
	return n.systemd.ReloadService("platform.nginx-public")
}

func (n *Nginx) InitConfig() error {
	domain := n.userConfig.GetDeviceDomain()
	configDir := n.systemConfig.ConfigDir()
	templateFile, err := os.ReadFile(path.Join(configDir, "nginx", "public.conf"))
	if err != nil {
		return err
	}
	tmpl := string(templateFile)
	tmpl = strings.ReplaceAll(tmpl, "{{ domain_regex }}", strings.ReplaceAll(domain, ".", "\\."))
	tmpl = strings.ReplaceAll(tmpl, "{{ domain }}", domain)
	nginxConfigDir := n.systemConfig.DataDir()
	nginxConfigFile := path.Join(nginxConfigDir, "nginx.conf")
	return os.WriteFile(nginxConfigFile, []byte(tmpl), 0644)
}

type customProxyTemplateData struct {
	Entries []customProxyServerEntry
}

type customProxyServerEntry struct {
	ServerName string
	Host       string
	Port       int
}

func (n *Nginx) InitCustomProxyConfig() error {
	return n.writeCustomProxyConfig()
}

func (n *Nginx) ReloadCustomProxy() error {
	err := n.writeCustomProxyConfig()
	if err != nil {
		return err
	}
	return n.systemd.ReloadService("platform.nginx-custom-proxy")
}

func (n *Nginx) writeCustomProxyConfig() error {
	domain := n.userConfig.GetDeviceDomain()
	configDir := n.systemConfig.ConfigDir()
	templateFile := path.Join(configDir, "nginx", "custom-proxy.conf")

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	entries, err := n.proxyConfig.Proxies()
	if err != nil {
		return err
	}

	serverEntries := make([]customProxyServerEntry, len(entries))
	for i, e := range entries {
		serverEntries[i] = customProxyServerEntry{
			ServerName: e.Name + "." + domain,
			Host:       e.Host,
			Port:       e.Port,
		}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, customProxyTemplateData{Entries: serverEntries})
	if err != nil {
		return err
	}

	nginxConfigDir := n.systemConfig.DataDir()
	nginxConfigFile := path.Join(nginxConfigDir, "custom-proxy.conf")
	return os.WriteFile(nginxConfigFile, buf.Bytes(), 0644)
}
