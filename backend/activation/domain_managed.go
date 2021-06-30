package activation

import (
	"fmt"
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/redirect"
	"log"
	"strings"
)

type ManagedActivateRequest struct {
	RedirectEmail    string `json:"redirect_email"`
	RedirectPassword string `json:"redirect_password"`
	Domain           string `json:"domain"`
	DeviceUsername   string `json:"device_username"`
	DevicePassword   string `json:"device_password"`
}

type ManagedPlatformUserConfig interface {
	SetRedirectEnabled(enabled bool)
	SetUserUpdateToken(userUpdateToken string)
	SetUserEmail(userEmail string)
	SetDomain(domain string)
	UpdateDomainToken(token string)
	GetRedirectDomain() string
}

type ManagedRedirect interface {
	Authenticate(email string, password string) (*redirect.User, error)
	Acquire(email string, password string, domain string) (*redirect.Domain, error)
	Reset(updateToken string) error
}

type Managed struct {
	internet connection.Checker
	config   ManagedPlatformUserConfig
	redirect ManagedRedirect
	device   DeviceActivation
}

func NewFree(internet connection.Checker, config ManagedPlatformUserConfig, redirect ManagedRedirect, device DeviceActivation) *Managed {
	return &Managed{
		internet: internet,
		config:   config,
		redirect: redirect,
		device:   device,
	}
}

func (f *Managed) Free(redirectEmail string, redirectPassword string, requestDomain string, deviceUsername string, devicePassword string) error {
	domain := fmt.Sprintf("%s.%s", strings.ToLower(requestDomain), f.config.GetRedirectDomain())
	return f.activate(redirectEmail, redirectPassword, domain, deviceUsername, devicePassword)
}

func (f *Managed) Premium(redirectEmail string, redirectPassword string, domain string, deviceUsername string, devicePassword string) error {
	return f.activate(redirectEmail, redirectPassword, domain, deviceUsername, devicePassword)
}

func (f *Managed) activate(redirectEmail string, redirectPassword string, domain string, deviceUsername string, devicePassword string) error {
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
