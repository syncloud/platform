package activation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
)

type CustomPlatformUserConfigStub struct {
	domain string
}

func (c *CustomPlatformUserConfigStub) SetRedirectEnabled(enabled bool) {
}

func (c *CustomPlatformUserConfigStub) SetUserEmail(userEmail string) {

}

func (c *CustomPlatformUserConfigStub) SetCustomDomain(domain string) {
	c.domain = domain
}

type CustorCertbotStub struct {
	attempted int
	generated int
	fail      bool
}

func (c *CustorCertbotStub) Generate() error {
	c.attempted += 1
	if c.fail {
		return fmt.Errorf("error")
	}
	c.generated++
	return nil
}

func TestManaged_ActivateCustom_GenerateFakeCertificate(t *testing.T) {
	logger := log.Default()

	cert := &CustorCertbotStub{}
	config := &CustomPlatformUserConfigStub{}
	managed := NewCustom(&InternetCheckerStub{}, config, &DeviceActivationStub{}, cert, logger)
	err := managed.Activate("example.com", "username", "password")
	assert.Nil(t, err)

	assert.Equal(t, 1, cert.generated)
}

func TestManaged_ActivateCustom_FixDomainName(t *testing.T) {
	logger := log.Default()

	cert := &CustorCertbotStub{}
	config := &CustomPlatformUserConfigStub{}
	managed := NewCustom(&InternetCheckerStub{}, config, &DeviceActivationStub{}, cert, logger)
	err := managed.Activate("ExaMple.com", "username", "password")
	assert.Nil(t, err)

	assert.Equal(t, "example.com", config.domain)
	assert.Equal(t, 1, cert.generated)
}
