package activation

import (
	"github.com/syncloud/platform/cert"
	"github.com/syncloud/platform/connection"
	"log"
	"strings"
)

type CustomActivateRequest struct {
	Domain         string `json:"domain"`
	DeviceUsername string `json:"device_username"`
	DevicePassword string `json:"device_password"`
}

type CustomPlatformUserConfig interface {
	SetRedirectEnabled(enabled bool)
	SetUserEmail(userEmail string)
	SetCustomDomain(domain string)
}

type CustomActivation interface {
	Activate(requestDomain string, deviceUsername string, devicePassword string) error
}

type Custom struct {
	internet connection.InternetChecker
	config   CustomPlatformUserConfig
	device   DeviceActivation
	cert     cert.Generator
}

func NewCustom(internet connection.InternetChecker, config CustomPlatformUserConfig, device DeviceActivation, cert cert.Generator) *Custom {
	return &Custom{
		internet: internet,
		config:   config,
		device:   device,
		cert:     cert,
	}
}

func (c *Custom) Activate(requestDomain string, deviceUsername string, devicePassword string) error {
	log.Printf("activate custom: %s, %s", requestDomain, deviceUsername)
	domain := strings.ToLower(requestDomain)

	err := c.internet.Check()
	if err != nil {
		return err
	}

	c.config.SetRedirectEnabled(false)
	c.config.SetCustomDomain(domain)
	name, email := ParseUsername(deviceUsername, domain)
	c.config.SetUserEmail(email)

	err = c.cert.Generate()
	if err != nil {
		return err
	}

	return c.device.ActivateDevice(deviceUsername, devicePassword, name, email)
}
