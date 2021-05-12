package activation

import (
	"fmt"
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/redirect"
	"log"
	"strings"
)

type FreePlatformUserConfig interface {
	SetRedirectEnabled(enabled bool)
	UpdateRedirect(domain string, apiUrl string)
	SetUserUpdateToken(userUpdateToken string)
	SetUserEmail(userEmail string)
}

type FreeRedirect interface {
	Authenticate(email string, password string) (*redirect.User, error)
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

func (f *Free) Activate(redirectEmail string, redirectPassword string, userDomain string, deviceUsername string, devicePassword string, mainDomain string) error {
	userDomainLower := strings.ToLower(userDomain)
	log.Printf("activate %s, %s", userDomainLower, deviceUsername)

	err := f.internet.Check()
	if err != nil {
		return err
	}
	user, err := f.prepareRedirect(redirectEmail, redirectPassword, mainDomain)
	if err != nil {
		return err
	}
	f.config.SetUserUpdateToken(user.UpdateToken)

	/*

	   name, email = parse_username(device_username, '{0}.{1}'.format(user_domain_lower, main_domain))

	   response_data = self.redirect_service.acquire(redirect_email, redirect_password, user_domain_lower)
	   self.user_platform_config.update_domain(response_data.user_domain, response_data.update_token)

	   self._activate_common(name, device_username, device_password, email)
	*/
	return fmt.Errorf("not implemented yet")
}

func (f *Free) prepareRedirect(redirectEmail string, redirectPassword string, mainDomain string) (*redirect.User, error) {
	redirectApiUrl := fmt.Sprintf("https://api.%s", mainDomain)

	log.Printf("prepare redirect %s, %s", redirectEmail, redirectApiUrl)
	f.config.SetRedirectEnabled(true)
	f.config.UpdateRedirect(mainDomain, redirectApiUrl)
	f.config.SetUserEmail(redirectEmail)

	user, err := f.redirect.Authenticate(redirectEmail, redirectPassword)
	if err != nil {
		return nil, err
	}
	return user, nil

}
