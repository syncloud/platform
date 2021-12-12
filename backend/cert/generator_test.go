package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"testing"
	"time"
)

type RedirectCertbotStub struct {
}

func (r RedirectCertbotStub) CertbotPresent(token, fqdn string, value ...string) error {
	//TODO implement me
	panic("implement me")
}

func (r RedirectCertbotStub) CertbotCleanUp(token, fqdn string) error {
	//TODO implement me
	panic("implement me")
}

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
	logConfig := zap.NewProductionConfig()
	logConfig.Encoding = "console"
	logConfig.EncoderConfig.TimeKey = ""
	logger, err := logConfig.Build()
	assert.Nil(t, err)
	now := time.Now()

	file := generateCertificate(now, Month-1*Day)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: file.Name(),
	}

	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{}
	fake := &FakeStub{}
	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err = generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.count)
}

func TestNotRegenerate_MoreThanAMonthBeforeExpiry(t *testing.T) {

	logConfig := zap.NewProductionConfig()
	logConfig.Encoding = "console"
	logConfig.EncoderConfig.TimeKey = ""
	logger, err := logConfig.Build()
	assert.Nil(t, err)
	now := time.Now()

	file := generateCertificate(now, Month+1*Day)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: file.Name(),
	}
	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{}
	fake := &FakeStub{}

	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err = generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.count)
}

func TestRegenerateFakeFallback(t *testing.T) {

	logger, err := zap.NewProduction()
	assert.Nil(t, err)
	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: "/unknown",
	}
	userConfig := &GeneratorUserConfigStub{activated: true}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}

	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err = generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 1, fake.count)
}

func TestRegenerateFake_IfDeviceIsNotActivated(t *testing.T) {

	logger, err := zap.NewProduction()
	assert.Nil(t, err)
	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: "/unknown",
	}
	userConfig := &GeneratorUserConfigStub{activated: false}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}

	generator := New(systemConfig, userConfig, provider, certbot, fake, logger)
	err = generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.attempt)
	assert.Equal(t, 0, certbot.count)
	assert.Equal(t, 1, fake.count)
}

func generateCertificate(now time.Time, duration time.Duration) *os.File {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             now,
		NotAfter:              now.Add(duration),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privateKey.Public(), privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	certFile, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(certFile.Name(), derBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return certFile
}
