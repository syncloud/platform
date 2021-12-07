package cert

import (
	"crypto/x509"
	"github.com/syncloud/platform/date"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	Log   = "/var/snap/platform/common/log/certbot.log"
	Day   = time.Hour * 24
	Month = Day * 30
)

type Generator interface {
	Generate() error
}

type CertificateGenerator struct {
	systemConfig GeneratorSystemConfig
	certbot      CertbotGenerator
	fake         FakeGenerator
	dateProvider date.Provider
	logger       *log.Logger
}

type GeneratorSystemConfig interface {
	SslCertificateFile() string
	SslKeyFile() string
}

func New(systemConfig GeneratorSystemConfig, dateProvider date.Provider, certbot CertbotGenerator, fake FakeGenerator) *CertificateGenerator {
	return &CertificateGenerator{
		systemConfig: systemConfig,
		certbot:      certbot,
		fake:         fake,
		dateProvider: dateProvider,
		logger:       log.Default(),
	}
}

func (g *CertificateGenerator) Start() {
	file, err := os.OpenFile(Log, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("unable to create certbot logger %v\n", err)
	}
	g.logger = log.New(file, "", log.LstdFlags)
}

func (g *CertificateGenerator) Generate() error {
	if !g.isExpired() {
		return nil
	}

	err := g.certbot.Generate()
	if err != nil {
		g.Log("unable to generate fake certificate: %v\n", err)
		return g.generateFake()
	}
	return nil
}

func (g *CertificateGenerator) isExpired() bool {

	certBytes, err := ioutil.ReadFile(g.systemConfig.SslCertificateFile())
	if err != nil {
		return true
	}

	certificate, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return true
	}

	now := g.dateProvider.Now()
	validFor := certificate.NotAfter.Sub(now)
	valid := validFor > Month
	realCert := certificate.Subject.String() != Subject
	if valid && realCert {
		g.Log("not regenerating real certificate, valid for days: %d\n", int(validFor.Hours()/24))
		return false
	}

	return true
}

func (g *CertificateGenerator) generateFake() error {
	err := g.fake.Generate()
	if err != nil {
		g.Log("unable to generate fake certificate: %v\n", err)
		return err
	}
	return nil
}

func (g *CertificateGenerator) Log(format string, v ...interface{}) {
	g.logger.Printf(format, v...)
}
