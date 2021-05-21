package activation

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/certificate"
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/nginx"
	"github.com/syncloud/platform/redirect"
	"log"
	"strings"
)

type FreeActivateRequest struct {
	RedirectEmail    string `json:"redirect_email"`
	RedirectPassword string `json:"redirect_password"`
	Domain           string `json:"user_domain"`
	DeviceUsername   string `json:"device_username"`
	DevicePassword   string `json:"device_password"`
}

type FreePlatformUserConfig interface {
	SetRedirectEnabled(enabled bool)
	UpdateRedirectDomain(domain string)
	UpdateRedirectApiUrl(apiUrl string)
	SetUserUpdateToken(userUpdateToken string)
	SetUserEmail(userEmail string)
	SetUserDomain(domain string)
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
	Acquire(email string, password string, userDomain string) (*redirect.Domain, error)
	Reset(updateToken string) error
}

type Free struct {
	internet             connection.Checker
	config               FreePlatformUserConfig
	redirect             FreeRedirect
	certificateGenerator *certificate.Generator
	auth                 *auth.LdapAuth
	nginx                *nginx.Nginx
}

func New(internet connection.Checker, config FreePlatformUserConfig, redirect FreeRedirect, certificateGenerator *certificate.Generator, auth *auth.LdapAuth, nginx *nginx.Nginx) *Free {
	return &Free{
		internet:             internet,
		config:               config,
		redirect:             redirect,
		certificateGenerator: certificateGenerator,
		auth:                 auth,
		nginx:                nginx,
	}
}

func (f *Free) Activate(redirectEmail string, redirectPassword string, userDomain string, deviceUsername string, devicePassword string) error {
	userDomainLower := strings.ToLower(userDomain)
	err := f.ActivateFreeDomain(redirectEmail, redirectPassword, userDomainLower)
	if err != nil {
		return err
	}
	return f.ActivateDevice(deviceUsername, devicePassword, userDomainLower)
}
func (f *Free) ActivateFreeDomain(redirectEmail string, redirectPassword string, userDomain string) error {
	log.Printf("activate: %s", userDomain)

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
	domain, err := f.redirect.Acquire(redirectEmail, redirectPassword, userDomain)
	if err != nil {
		return err
	}

	f.config.SetUserDomain(domain.UserDomain)
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

	log.Println("activating ldap")
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
	//self.event_trigger.trigger_app_event_domain()
	return fmt.Errorf("not implemented yet")
}

func ParseUsername(username string, domain string) (string, string) {
	if strings.Contains(username, "@") {
		return strings.Split(username, "@")[0], username
	}
	email := fmt.Sprintf("%s@%s", username, domain)
	return username, email
}
