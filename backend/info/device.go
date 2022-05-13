package info

import (
	"fmt"
)

func ConstructUrl(port *int, domain string, app string) string {
	externalPort := ""
	if port != nil && *port != 80 && *port != 443 {
		externalPort = fmt.Sprintf(":%d", *port)
	}
	return fmt.Sprintf("https://%s.%s%s", app, domain, externalPort)
}

type DeviceUserConfig interface {
	GetPublicPort() *int
	GetDeviceDomain() string
}

type Device struct {
	userConfig DeviceUserConfig
}

func New(userConfig DeviceUserConfig) *Device {
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
