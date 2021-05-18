package activation

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/syncloud/platform/connection"
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
	UpdateUserDomain(domain string)
	UpdateDomainToken(token string)
	GetRedirectDomain() string
	SetActivated()
	SetWebSecretKey(key string)
	GetDomainUpdateToken() *string
}

type FreeRedirect interface {
	Authenticate(email string, password string) (*redirect.User, error)
	Acquire(email string, password string, userDomain string) (*redirect.Domain, error)
	Reset(updateToken string) error
}

type Free struct {
	internet connection.Checker
	config   FreePlatformUserConfig
	redirect FreeRedirect
}

func New(internet connection.Checker, config FreePlatformUserConfig, redirect FreeRedirect) *Free {
	return &Free{
		internet: internet,
		config:   config,
		redirect: redirect,
	}
}

func (f *Free) Activate(redirectEmail string, redirectPassword string, userDomain string, deviceUsername string, devicePassword string) error {
	userDomainLower := strings.ToLower(userDomain)
	log.Printf("activate %s, %s", userDomainLower, deviceUsername)

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
	mainDomain := f.config.GetRedirectDomain()
	name, email := ParseUsername(deviceUsername, fmt.Sprintf("%s.%s", userDomainLower, mainDomain))
	domain, err := f.redirect.Acquire(redirectEmail, redirectPassword, userDomainLower)
	if err != nil {
		return err
	}

	f.config.UpdateUserDomain(domain.UserDomain)
	f.config.UpdateDomainToken(domain.UpdateToken)
	err = f.redirect.Reset(domain.UpdateToken)
	if err != nil {
		return err
	}

	return f.ActivateCommon(name, deviceUsername, devicePassword, email)
}

func (f *Free) ActivateCommon(name string, username string, password string, email string) error {
	f.resetAccess()

	log.Println("activating ldap")
	f.config.SetWebSecretKey(uuid.New().String())

	//self.tls.generate_self_signed_certificate()

	//self.auth.reset(name, device_username, device_password, email)

	//self.nginx.init_config()
	//self.nginx.reload_public()

	f.config.SetActivated()

	log.Println("activation completed")

	return fmt.Errorf("not implemented yet")
}

func (f *Free) resetAccess() {
	log.Println("reset access")
	//f.config.update_device_access(False, False, None, 0, 0)
	//self.event_trigger.trigger_app_event_domain()
}

func ParseUsername(username string, domain string) (string, string) {
	if strings.Contains(username, "@") {
		return strings.Split(username, "@")[0], username
	}
	email := fmt.Sprintf("%s@%s", username, domain)
	return username, email
}
