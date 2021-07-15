package activation

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/redirect"
	"testing"
)

type InternetCheckerStub struct{}

func (i *InternetCheckerStub) Check() error {
	return nil
}

type FreePlatformUserConfigStub struct {
}

func (f *FreePlatformUserConfigStub) SetRedirectEnabled(enabled bool) {

}

func (f *FreePlatformUserConfigStub) SetUserUpdateToken(userUpdateToken string) {
}

func (f *FreePlatformUserConfigStub) SetUserEmail(userEmail string) {
}

func (f *FreePlatformUserConfigStub) SetDomain(domain string) {
}

func (f *FreePlatformUserConfigStub) UpdateDomainToken(token string) {
}

func (f *FreePlatformUserConfigStub) GetRedirectDomain() string {
	return "syncloud.it"
}

type FreeRedirectStub struct {
	email    string
	password string
	domain   string
}

func (f *FreeRedirectStub) Authenticate(email string, password string) (*redirect.User, error) {
	return &redirect.User{UpdateToken: "user_token"}, nil
}

func (f *FreeRedirectStub) Acquire(email string, password string, domain string) (*redirect.Domain, error) {
	f.email = email
	f.password = password
	f.domain = domain
	return &redirect.Domain{
		Name:        domain,
		UpdateToken: "domain_token",
	}, nil
}

func (f *FreeRedirectStub) Reset(updateToken string) error {
	return nil
}

type DeviceActivationStub struct {
}

func (d *DeviceActivationStub) ActivateDevice(username string, password string, name string, email string) error {
	return nil
}

func TestFree_Activate(t *testing.T) {
	freeRedirect := &FreeRedirectStub{}
	free := NewFree(&InternetCheckerStub{}, &FreePlatformUserConfigStub{}, freeRedirect, &DeviceActivationStub{})
	err := free.Activate("mail", "password", "test.syncloud.it", "username", "password")
	assert.Nil(t, err)

	assert.Equal(t, "test.syncloud.it", freeRedirect.domain)
}
