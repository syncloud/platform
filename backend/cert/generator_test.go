package cert

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

type GeneratorUserConfigStub struct {
	activated bool
	domain    string
}

func (g *GeneratorUserConfigStub) IsActivated() bool {
	return g.activated
}

func (g *GeneratorUserConfigStub) GetDeviceDomain() string {
	return g.domain
}

type GeneratorSystemConfigStub struct {
	certFile   string
	keyFile    string
	caSertFile string
	caKeyFile  string
}

func (g GeneratorSystemConfigStub) SslCertificateFile() string {
	return g.certFile
}

func (g GeneratorSystemConfigStub) SslKeyFile() string {
	return g.keyFile
}

func (g GeneratorSystemConfigStub) SslCaCertificateFile() string {
	return g.caSertFile
}

func (g GeneratorSystemConfigStub) SslCaKeyFile() string {
	return g.caKeyFile
}

type ProviderStub struct {
	now time.Time
}

func (p ProviderStub) Now() time.Time {
	return p.now
}

type CertbotStub struct {
	attempt int
	count   int
	fail    bool
}

func (c *CertbotStub) Generate() error {
	c.attempt++
	if c.fail {
		return fmt.Errorf("certbot fail")
	}
	c.count++
	return nil
}

type FakeStub struct {
	count int
}

func (f *FakeStub) Generate() error {
	f.count++
	return nil
}

type GeneratorNginxStub struct {
	reloadPublic int
}

func (n *GeneratorNginxStub) ReloadPublic() error {
	n.reloadPublic++
	return nil
}

type TriggerStub struct {
	called int
}

func (t *TriggerStub) RunCertificateChangeEvent() error {
	t.called++
	return nil
}

func TestRegenerate_LessThanAMonthBeforeExpiry(t *testing.T) {
	logger := log.Default()
	now := time.Now()

	file := generateCertificate(now, Month-1*Day, true)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		certFile: file.Name(),
	}

	userConfig := &GeneratorUserConfigStub{activated: true, domain: "domain"}
	certbot := &CertbotStub{}
	fake := &FakeStub{}
	nginx := &GeneratorNginxStub{}
	trigger := &TriggerStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, nginx, trigger, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.count)
	assert.Equal(t, 1, nginx.reloadPublic)
	assert.Equal(t, 1, trigger.called)

}

func TestRegenerate_WrongDomain(t *testing.T) {
	logger := log.Default()
	now := time.Now()

	file := generateCertificate(now, Month+1*Day, true)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		certFile: file.Name(),
	}

	userConfig := &GeneratorUserConfigStub{activated: true, domain: "new.domain"}
	certbot := &CertbotStub{}
	fake := &FakeStub{}
	nginx := &GeneratorNginxStub{}
	trigger := &TriggerStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, nginx, trigger, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.count)
	assert.Equal(t, 1, nginx.reloadPublic)
	assert.Equal(t, 1, trigger.called)

}

func TestNotRegenerate_MoreThanAMonthBeforeExpiry(t *testing.T) {

	logger := log.Default()
	now := time.Now()

	file := generateCertificate(now, Month+1*Day, true)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		certFile: file.Name(),
	}
	userConfig := &GeneratorUserConfigStub{activated: true, domain: "domain"}
	certbot := &CertbotStub{}
	fake := &FakeStub{}
	nginx := &GeneratorNginxStub{}
	trigger := &TriggerStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, nginx, trigger, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 0, nginx.reloadPublic)
	assert.Equal(t, 0, trigger.called)
}

func TestRegenerateFakeFallback(t *testing.T) {

	logger := log.Default()
	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		certFile: "/unknown",
	}
	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}
	nginx := &GeneratorNginxStub{}
	trigger := &TriggerStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, nginx, trigger, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 1, fake.count)
	assert.Equal(t, 1, nginx.reloadPublic)
	assert.Equal(t, 1, trigger.called)
}

func TestNotGenerateFakeIfValid(t *testing.T) {

	logger := log.Default()
	now := time.Now()
	provider := &ProviderStub{now: now}

	systemConfig := &GeneratorSystemConfigStub{
		certFile: generateCertificate(now, Month+1*Day, false).Name(),
	}
	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}
	nginx := &GeneratorNginxStub{}
	trigger := &TriggerStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, nginx, trigger, logger)
	err := generator.Generate()
	assert.Nil(t, err)

	assert.Equal(t, 1, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 0, fake.count)
	assert.Equal(t, 0, nginx.reloadPublic)
	assert.Equal(t, 0, trigger.called)
}

func TestRegenerateFake_IfDeviceIsNotActivated(t *testing.T) {

	logger := log.Default()

	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		certFile: "/unknown",
	}
	userConfig := &GeneratorUserConfigStub{activated: false}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}
	nginx := &GeneratorNginxStub{}
	trigger := &TriggerStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, nginx, trigger, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 1, fake.count)
	assert.Equal(t, 0, nginx.reloadPublic)
	assert.Equal(t, 0, trigger.called)
}

func TestNotGenerateFake_IfDeviceIsNotActivatedButCertIsValid(t *testing.T) {

	logger := log.Default()

	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		certFile: generateCertificate(now, Month+1*Day, false).Name(),
	}
	userConfig := &GeneratorUserConfigStub{activated: false}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}
	nginx := &GeneratorNginxStub{}
	trigger := &TriggerStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, nginx, trigger, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 0, fake.count)
}

func generateCertificate(now time.Time, duration time.Duration, real bool) *os.File {

	subjectOrganization := SubjectOrganization
	if real {
		subjectOrganization = "Real"
	}

	certFile, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	keyFile, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	caCertFile, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	caKeyFile, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	fake := NewFake(
		&GeneratorSystemConfigStub{
			certFile:   certFile.Name(),
			keyFile:    keyFile.Name(),
			caSertFile: caCertFile.Name(),
			caKeyFile:  caKeyFile.Name(),
		},
		&GeneratorUserConfigStub{
			domain: "domain",
		},
		&ProviderStub{now: now},
		subjectOrganization,
		duration,
		log.Default(),
	)
	err = fake.Generate()
	if err != nil {
		panic(err)
	}

	return certFile
}
