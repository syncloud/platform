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

type ManagedPlatformUserConfigStub struct {
}

func (f *ManagedPlatformUserConfigStub) SetRedirectEnabled(enabled bool) {

}

func (f *ManagedPlatformUserConfigStub) SetUserUpdateToken(userUpdateToken string) {
}

func (f *ManagedPlatformUserConfigStub) SetUserEmail(userEmail string) {
}

func (f *ManagedPlatformUserConfigStub) SetDomain(domain string) {
}

func (f *ManagedPlatformUserConfigStub) UpdateDomainToken(token string) {
}

func (f *ManagedPlatformUserConfigStub) GetRedirectDomain() string {
	return "syncloud.it"
}

type ManagedRedirectStub struct {
	email    string
	password string
	domain   string
}

func (f *ManagedRedirectStub) Authenticate(email string, password string) (*redirect.User, error) {
	return &redirect.User{UpdateToken: "user_token"}, nil
}

func (f *ManagedRedirectStub) Acquire(email string, password string, domain string) (*redirect.Domain, error) {
	f.email = email
	f.password = password
	f.domain = domain
	return &redirect.Domain{
		Name:        domain,
		UpdateToken: "domain_token",
	}, nil
}

func (f *ManagedRedirectStub) Reset(updateToken string) error {
	return nil
}

type DeviceActivationStub struct {
}

type ManagedCertbotStub struct {
}

func (c *ManagedCertbotStub) Generate(email, domain, token string) error {
	return nil
}

func (d *DeviceActivationStub) ActivateDevice(username string, password string, name string, email string) error {
	return nil
}

func TestManaged_Activate(t *testing.T) {
	managedRedirect := &ManagedRedirectStub{}
	managed := NewManaged(&InternetCheckerStub{}, &ManagedPlatformUserConfigStub{}, managedRedirect, &DeviceActivationStub{}, &ManagedCertbotStub{})
	err := managed.Activate("mail", "password", "test.syncloud.it", "username", "password")
	assert.Nil(t, err)

	assert.Equal(t, "test.syncloud.it", managedRedirect.domain)
}
