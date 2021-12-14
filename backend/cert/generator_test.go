package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"
)

type GeneratorUserConfigStub struct {
	activated bool
}

func (g *GeneratorUserConfigStub) IsActivated() bool {
	return g.activated
}

type GeneratorSystemConfigStub struct {
	sslCertificateFile string
}

func (g GeneratorSystemConfigStub) SslCertificateFile() string {
	return g.sslCertificateFile
}

func (g GeneratorSystemConfigStub) SslKeyFile() string {
	//TODO implement me
	panic("implement me")
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

func TestRegenerate_LessThanAMonthBeforeExpiry(t *testing.T) {
	logger := log.Default()
	now := time.Now()

	file := generateCertificate(now, Month-1*Day, true)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: file.Name(),
	}

	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{}
	fake := &FakeStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.count)
}

func TestNotRegenerate_MoreThanAMonthBeforeExpiry(t *testing.T) {

	logger := log.Default()
	now := time.Now()

	file := generateCertificate(now, Month+1*Day, true)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: file.Name(),
	}
	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{}
	fake := &FakeStub{}

	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.count)
}

func TestRegenerateFakeFallback(t *testing.T) {

	logger := log.Default()
	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: "/unknown",
	}
	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}

	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 1, fake.count)
}

func TestNotGenerateFakeIfValid(t *testing.T) {

	logger := log.Default()
	now := time.Now()
	provider := &ProviderStub{now: now}

	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: generateCertificate(now, Month+1*Day, false).Name(),
	}
	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}

	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err := generator.Generate()
	assert.Nil(t, err)

	assert.Equal(t, 1, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 0, fake.count)
}

func TestRegenerateFake_IfDeviceIsNotActivated(t *testing.T) {

	logger := log.Default()

	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: "/unknown",
	}
	userConfig := &GeneratorUserConfigStub{activated: false}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}

	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 1, fake.count)
}

func TestNotGenerateFake_IfDeviceIsNotActivatedButCertIsValid(t *testing.T) {

	logger := log.Default()

	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: generateCertificate(now, Month+1*Day, false).Name(),
	}
	userConfig := &GeneratorUserConfigStub{activated: false}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}

	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 0, fake.count)
}

func generateCertificate(now time.Time, duration time.Duration, real bool) *os.File {

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		panic(err)
	}

	subject := pkix.Name{
		Organization: []string{"Acme Co"},
	}

	if !real {
		subject = pkix.Name{
			Country:      []string{SubjectCountry},
			Province:     []string{SubjectProvince},
			Locality:     []string{SubjectLocality},
			Organization: []string{SubjectOrganization},
			CommonName:   SubjectCommonName,
		}
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               subject,
		NotBefore:             now,
		NotAfter:              now.Add(duration),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privateKey.Public(), privateKey)
	if err != nil {
		panic(err)
	}
	certFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(certFile.Name(), derBytes, 0644)
	if err != nil {
		panic(err)
	}

	return certFile
}
