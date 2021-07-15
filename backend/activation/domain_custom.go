package activation

import (
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
	SetUserUpdateToken(userUpdateToken string)
	SetUserEmail(userEmail string)
	SetDomain(domain string)
	UpdateDomainToken(token string)
	GetRedirectDomain() string
	SetCustomDomain(domain string)
}

type CustomActivation interface {
	Activate(requestDomain string, deviceUsername string, devicePassword string) error
}

type Custom struct {
	internet connection.Checker
	config   CustomPlatformUserConfig
	redirect ManagedRedirect
	device   *Device
}

func NewCustom(internet connection.Checker, config CustomPlatformUserConfig, redirect ManagedRedirect, device *Device) *Custom {
	return &Custom{
		internet: internet,
		config:   config,
		redirect: redirect,
		device:   device,
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
	return c.device.ActivateDevice(deviceUsername, devicePassword, name, email)
}
