package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/rest/model"
)

type fakeMailRelayUserConfig struct {
	DeviceUserConfig
	domain       string
	relayEnabled bool
	token        *string
}

func (f *fakeMailRelayUserConfig) GetDeviceDomain() string       { return f.domain }
func (f *fakeMailRelayUserConfig) IsMailRelayEnabled() bool      { return f.relayEnabled }
func (f *fakeMailRelayUserConfig) GetDomainUpdateToken() *string { return f.token }

type fakeMailRelayRedirect struct {
	DeviceRedirect
	domain string
}

func (f *fakeMailRelayRedirect) Domain() string { return f.domain }

func mailRelay(t *testing.T, userConfig DeviceUserConfig, redirect DeviceRedirect) *model.MailRelayCredentials {
	t.Helper()
	api := &Api{userConfig: userConfig, redirect: redirect}
	response, err := api.MailRelay(nil)
	assert.NoError(t, err)
	credentials, ok := response.(*model.MailRelayCredentials)
	assert.True(t, ok)
	return credentials
}

func TestMailRelay_Enabled(t *testing.T) {
	token := "the-token"
	credentials := mailRelay(t,
		&fakeMailRelayUserConfig{domain: "device.syncloud.it", relayEnabled: true, token: &token},
		&fakeMailRelayRedirect{domain: "syncloud.it"})

	assert.True(t, credentials.Enabled)
	assert.Equal(t, "mail-relay.syncloud.it", credentials.Host)
	assert.Equal(t, 587, credentials.Port)
	assert.Equal(t, "device.syncloud.it", credentials.Login)
	assert.Equal(t, token, credentials.Password)
}

func TestMailRelay_DisabledByToggle(t *testing.T) {
	token := "the-token"
	credentials := mailRelay(t,
		&fakeMailRelayUserConfig{domain: "device.syncloud.it", relayEnabled: false, token: &token},
		&fakeMailRelayRedirect{domain: "syncloud.it"})

	assert.False(t, credentials.Enabled)
	assert.Empty(t, credentials.Host)
	assert.Empty(t, credentials.Password)
}

func TestMailRelay_DisabledWithoutToken(t *testing.T) {
	credentials := mailRelay(t,
		&fakeMailRelayUserConfig{domain: "device.syncloud.it", relayEnabled: true, token: nil},
		&fakeMailRelayRedirect{domain: "syncloud.it"})

	assert.False(t, credentials.Enabled)
	assert.Empty(t, credentials.Password)
}
