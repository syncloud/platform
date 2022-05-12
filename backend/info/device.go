package info

import (
	"fmt"
	"github.com/syncloud/platform/config"
)

func ConstructUrl(port *int, domain string, app string) string {
	externalPort := ""
	if port != nil && *port != 80 && *port != 443 {
		externalPort = fmt.Sprintf(":%d", port)
	}
	return fmt.Sprintf("https://%s.%s%s}", app, domain, externalPort)
}

type Device struct {
	userConfig *config.UserConfig
}

func New(userConfig *config.UserConfig) *Device {
	return &Device{
		userConfig: userConfig,
	}
}

func (d *Device) AppDomain(app string) string {
	return fmt.Sprintf("%s.%s", app, d.userConfig.GetDeviceDomain())
}

func (d *Device) Url(app string) string {
	port := d.userConfig.GetPublicPort()
	domain := d.userConfig.GetDeviceDomain()
	return ConstructUrl(port, domain, app)
}
