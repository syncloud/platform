package activation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type CustomPlatformUserConfigStub struct {
}

func (c CustomPlatformUserConfigStub) SetRedirectEnabled(enabled bool) {
}

func (c CustomPlatformUserConfigStub) SetUserEmail(userEmail string) {

}

func (c CustomPlatformUserConfigStub) SetCustomDomain(domain string) {
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
	cert := &CustorCertbotStub{}
	config := &CustomPlatformUserConfigStub{}
	managed := NewCustom(&InternetCheckerStub{}, config, &DeviceActivationStub{}, cert)
	err := managed.Activate("example.com", "username", "password")
	assert.Nil(t, err)

	assert.Equal(t, 1, cert.generated)
}
