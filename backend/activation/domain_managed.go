package activation

import (
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/redirect"
	"log"
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

type ManagedActivation interface {
	Activate(redirectEmail string, redirectPassword string, requestDomain string, deviceUsername string, devicePassword string) error
}

type ManagedCertbot interface {
	Generate(email string, domain string, token string) error
}

type Managed struct {
	internet connection.InternetChecker
	config   ManagedPlatformUserConfig
	redirect ManagedRedirect
	device   DeviceActivation
	certbot  ManagedCertbot
}

func NewManaged(internet connection.InternetChecker, config ManagedPlatformUserConfig, redirect ManagedRedirect, device DeviceActivation, certbot ManagedCertbot) *Managed {
	return &Managed{
		internet: internet,
		config:   config,
		redirect: redirect,
		device:   device,
		certbot:  certbot,
	}
}

func (f *Managed) Activate(redirectEmail string, redirectPassword string, domainName string, deviceUsername string, devicePassword string) error {
	log.Printf("activate: %s", domainName)

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
	domain, err := f.redirect.Acquire(redirectEmail, redirectPassword, domainName)
	if err != nil {
		return err
	}

	f.config.SetDomain(domain.Name)
	f.config.UpdateDomainToken(domain.UpdateToken)
	err = f.redirect.Reset(domain.UpdateToken)
	if err != nil {
		return err
	}

	name, email := ParseUsername(deviceUsername, domain.Name)

	err = f.certbot.Generate(email, domain.Name, domain.UpdateToken)
	if err != nil {
		return err
	}

	return f.device.ActivateDevice(deviceUsername, devicePassword, name, email)
}
