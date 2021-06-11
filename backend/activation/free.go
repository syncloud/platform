package activation

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/certificate"
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/nginx"
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
	UpdateRedirectDomain(domain string)
	UpdateRedirectApiUrl(apiUrl string)
	SetUserUpdateToken(userUpdateToken string)
	SetUserEmail(userEmail string)
	SetDomain(domain string)
	UpdateDomainToken(token string)
	GetRedirectDomain() string
	SetActivated()
	SetWebSecretKey(key string)
	GetDomainUpdateToken() *string
	SetExternalAccess(enabled bool)
	SetUpnp(enabled bool)
	SetPublicIp(publicIp string)
	SetManualCertificatePort(manualCertificatePort int)
	SetManualAccessPort(manualAccessPort int)
	DeletePublicIp()
}

type FreeRedirect interface {
	Authenticate(email string, password string) (*redirect.User, error)
	Acquire(email string, password string, domain string) (*redirect.Domain, error)
	Reset(updateToken string) error
}

type Free struct {
	internet             connection.Checker
	config               FreePlatformUserConfig
	redirect             FreeRedirect
	certificateGenerator *certificate.Generator
	auth                 *auth.Service
	nginx                *nginx.Nginx
	trigger              *event.Trigger
}

func New(internet connection.Checker, config FreePlatformUserConfig, redirect FreeRedirect, certificateGenerator *certificate.Generator, auth *auth.Service, nginx *nginx.Nginx, trigger *event.Trigger) *Free {
	return &Free{
		internet:             internet,
		config:               config,
		redirect:             redirect,
		certificateGenerator: certificateGenerator,
		auth:                 auth,
		nginx:                nginx,
		trigger:              trigger,
	}
}

func (f *Free) Activate(redirectEmail string, redirectPassword string, requestDomain string, deviceUsername string, devicePassword string) error {
	domain := fmt.Sprintf("%s.%s",strings.ToLower(requestDomain), f.config.GetRedirectDomain())
  
	err := f.ActivateFreeDomain(redirectEmail, redirectPassword, domain)
	if err != nil {
		return err
	}
	return f.ActivateDevice(deviceUsername, devicePassword, domain)
}
func (f *Free) ActivateFreeDomain(redirectEmail string, redirectPassword string, requestDomain string) error {
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

func (f *Free) ActivateDevice(username string, password string, userDomain string) error {
	mainDomain := f.config.GetRedirectDomain()
	name, email := ParseUsername(username, fmt.Sprintf("%s.%s", userDomain, mainDomain))
	err := f.resetAccess()
	if err != nil {
		return err
	}

	f.config.SetWebSecretKey(uuid.New().String())

	err = f.certificateGenerator.GenerateSelfSigned()
	if err != nil {
		return err
	}

	err = f.auth.Reset(name, username, password, email)
	if err != nil {
		return err
	}

	err = f.nginx.InitConfig()
	if err != nil {
		return err
	}
	err = f.nginx.ReloadPublic()
	if err != nil {
		return err
	}

	f.config.SetActivated()

	log.Println("activation completed")

	return nil
}

func (f *Free) resetAccess() error {
	log.Println("reset access")
	f.config.SetUpnp(false)
	f.config.SetExternalAccess(false)
	f.config.DeletePublicIp()
	f.config.SetManualCertificatePort(0)
	f.config.SetManualAccessPort(0)
	return f.trigger.RunAccessChangeEvent()
}

func ParseUsername(username string, domain string) (string, string) {
	if strings.Contains(username, "@") {
		return strings.Split(username, "@")[0], username
	}
	email := fmt.Sprintf("%s@%s", username, domain)
	return username, email
}
