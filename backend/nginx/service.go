package nginx

import (
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type Systemd interface {
	ReloadService(service string) error
}

type SystemConfig interface {
	ConfigDir() (*string, error)
	NginxConfigDir() (*string, error)
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

	configDir, err := n.systemConfig.ConfigDir()
	if err != nil {
		return err
	}

	templateFile, err := ioutil.ReadFile(path.Join(*configDir, "nginx", "public.conf"))
	if err != nil {
		return err
	}

	template := string(templateFile)
	template = strings.ReplaceAll(template, "{{ user_domain }}", strings.ReplaceAll(domain, ".", "\\."))
	log.Printf("nginx config: %s", template)
	nginxConfigDir, err := n.systemConfig.NginxConfigDir()
	if err != nil {
		return err
	}
	nginxConfigFile := path.Join(*nginxConfigDir, "nginx.conf")
	log.Printf("nginx config file: %s", nginxConfigFile)
	err = ioutil.WriteFile(nginxConfigFile, []byte(template), 0644)
	return err
}
