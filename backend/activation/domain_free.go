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
	log.Printf("activate: %s", domain)

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
	domainResponse, err := f.redirect.Acquire(redirectEmail, redirectPassword, domain)
	if err != nil {
		return err
	}

	f.config.SetDomain(domainResponse.Name)
	f.config.UpdateDomainToken(domainResponse.UpdateToken)
	err = f.redirect.Reset(domainResponse.UpdateToken)
	if err != nil {
		return err
	}

	name, email := ParseUsername(deviceUsername, domain)
	return f.device.ActivateDevice(deviceUsername, devicePassword, name, email)
}
