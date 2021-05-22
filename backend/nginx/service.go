package nginx

import (
	"fmt"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/systemd"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type Nginx struct {
	systemd      *systemd.Control
	systemConfig *config.SystemConfig
	userConfig   *config.UserConfig
}

func New(systemd *systemd.Control, systemConfig *config.SystemConfig, userConfig *config.UserConfig) *Nginx {
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
	if domain == nil {
		return fmt.Errorf("device domain is not set")
	}

	configDir, err := n.systemConfig.ConfigDir()
	if err != nil {
		return err
	}

	templateFile, err := ioutil.ReadFile(path.Join(*configDir, "nginx", "public.conf"))
	if err != nil {
		return err
	}

	template := string(templateFile)
	template = strings.ReplaceAll(template, "${user_domain}", *domain)
	log.Printf("nginx config: %s", template)
	nginxConfigDir, err := n.systemConfig.NginxConfigDir()
	if err != nil {
		return err
	}
	nginxConfigFile := path.Join(*nginxConfigDir, "nginx.conf")
	log.Printf("nginx config file: %s", nginxConfigFile)
	err = ioutil.WriteFile(nginxConfigFile, []byte(template), 644)
	return err
}
