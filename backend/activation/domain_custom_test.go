package activation

import (
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

func TestManaged_ActivateCustom_GenerateFakeCertificate(t *testing.T) {
	fakeCert := &FakeCertbotStub{}
	config := &CustomPlatformUserConfigStub{}
	managed := NewCustom(&InternetCheckerStub{}, config, &DeviceActivationStub{}, fakeCert)
	err := managed.Activate("example.com", "username", "password")
	assert.Nil(t, err)

	assert.Equal(t, 1, fakeCert.generated)
}
