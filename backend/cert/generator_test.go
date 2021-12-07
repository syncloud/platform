package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	domain         string
	redirectDomain string
}

func (g GeneratorUserConfigStub) GetDomain() *string {
	return &g.domain
}

func (g GeneratorUserConfigStub) GetRedirectDomain() string {
	return g.redirectDomain
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

	now := time.Now()

	file := generateCertificate(now, Month-1*Day)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: file.Name(),
	}

	certbot := &CertbotStub{}
	fake := &FakeStub{}
	generator := New(systemConfig, provider, certbot, fake)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.count)
}

func TestNotRegenerate_MoreThanAMonthBeforeExpiry(t *testing.T) {

	now := time.Now()

	file := generateCertificate(now, Month+1*Day)
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: file.Name(),
	}
	certbot := &CertbotStub{}
	fake := &FakeStub{}

	generator := New(systemConfig, provider, certbot, fake)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.count)
}

func TestRegenerateFakeFallback(t *testing.T) {

	now := time.Now()
	provider := &ProviderStub{now: now}
	systemConfig := &GeneratorSystemConfigStub{
		sslCertificateFile: "/unknown",
	}
	certbot := &CertbotStub{fail: true}
	fake := &FakeStub{}

	generator := New(systemConfig, provider, certbot, fake)
	err := generator.Generate()
	assert.Nil(t, err)
	assert.Equal(t, 1, certbot.attempt)
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
