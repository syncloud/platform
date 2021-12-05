package real

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
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
}

func (g GeneratorUserConfigStub) IsCertbotStaging() bool {
	//TODO implement me
	panic("implement me")
}

func (g GeneratorUserConfigStub) GetUserEmail() *string {
	//TODO implement me
	panic("implement me")
}

func (g GeneratorUserConfigStub) GetDomain() *string {
	//TODO implement me
	panic("implement me")
}

func (g GeneratorUserConfigStub) GetDomainUpdateToken() *string {
	//TODO implement me
	panic("implement me")
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
	count int
}

func (c *CertbotStub) Generate() error {
	c.count++
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
	generator := New(systemConfig, provider, certbot)
	err := generator.RegenerateIfNeeded()
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
	generator := New(systemConfig, provider, certbot)
	err := generator.RegenerateIfNeeded()
	assert.Nil(t, err)
	assert.Equal(t, 0, certbot.count)
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
