package activation

import (
	"fmt"
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/redirect"
	"log"
	"strings"
)

type FreeActivateRequest struct {
	RedirectEmail    string `json:"redirect_email"`
	RedirectPassword string `json:"redirect_password"`
	Domain           string `json:"domain"`
	DeviceUsername   string `json:"device_username"`
	DevicePassword   string `json:"device_password"`
}

type FreePlatformUserConfig interface {
	SetRedirectEnabled(enabled bool)
	SetUserUpdateToken(userUpdateToken string)
	SetUserEmail(userEmail string)
	SetDomain(domain string)
	UpdateDomainToken(token string)
	GetRedirectDomain() string
}

type FreeRedirect interface {
	Authenticate(email string, password string) (*redirect.User, error)
	Acquire(email string, password string, domain string) (*redirect.Domain, error)
	Reset(updateToken string) error
}

type Free struct {
	internet connection.Checker
	config   FreePlatformUserConfig
	redirect FreeRedirect
	device   *Device
}

func NewFree(internet connection.Checker, config FreePlatformUserConfig, redirect FreeRedirect, device *Device) *Free {
	return &Free{
		internet: internet,
		config:   config,
		redirect: redirect,
		device:   device,
	}
}

func (f *Free) Activate(redirectEmail string, redirectPassword string, requestDomain string, deviceUsername string, devicePassword string) error {
	domain := fmt.Sprintf("%s.%s", strings.ToLower(requestDomain), f.config.GetRedirectDomain())
	err := f.activateFreeDomain(redirectEmail, redirectPassword, domain)
	if err != nil {
		return err
	}
	name, email := ParseUsername(deviceUsername, domain)
	return f.device.ActivateDevice(deviceUsername, devicePassword, name, email)
}

func (f *Free) activateFreeDomain(redirectEmail string, redirectPassword string, requestDomain string) error {
	log.Printf("activate: %s", requestDomain)

	err := f.internet.Check()
	if err != nil {
		return err
	}

	f.config.SetRedirectEnabled(true)
	f.config.SetUserEmail(redirectEmail)
	user, err := f.redirect.Authenticate(redirectEmail, redirectPassword)
	if err != nil {
		return err
	}

	f.config.SetUserUpdateToken(user.UpdateToken)
	domain, err := f.redirect.Acquire(redirectEmail, redirectPassword, requestDomain)
	if err != nil {
		return err
	}

	f.config.SetDomain(domain.Name)
	f.config.UpdateDomainToken(domain.UpdateToken)
	return f.redirect.Reset(domain.UpdateToken)
}
