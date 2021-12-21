package nginx

import (
	"io/ioutil"
	"path"
	"strings"
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

type Nginx struct {
	systemd      Systemd
	systemConfig SystemConfig
	userConfig   UserConfig
}

func New(systemd Systemd, systemConfig SystemConfig, userConfig UserConfig) *Nginx {
	return &Nginx{
		systemd:      systemd,
		userConfig:   userConfig,
		systemConfig: systemConfig,
	}
}

func (n *Nginx) ReloadPublic() error {
	return n.systemd.ReloadService("platform.nginx-public")
}

func (n *Nginx) InitConfig() error {
	domain := n.userConfig.GetDeviceDomain()
	configDir := n.systemConfig.ConfigDir()
	templateFile, err := ioutil.ReadFile(path.Join(configDir, "nginx", "public.conf"))
	if err != nil {
		return err
	}
	template := string(templateFile)
	template = strings.ReplaceAll(template, "{{ domain }}", strings.ReplaceAll(domain, ".", "\\."))
	nginxConfigDir := n.systemConfig.DataDir()
	nginxConfigFile := path.Join(nginxConfigDir, "nginx.conf")
	err = ioutil.WriteFile(nginxConfigFile, []byte(template), 0644)
	return err
}
