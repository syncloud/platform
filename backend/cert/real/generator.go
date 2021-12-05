package real

import (
	"crypto"
	"crypto/x509"
	"github.com/go-acme/lego/v4/registration"
	"github.com/syncloud/platform/cert"
	"github.com/syncloud/platform/cert/fake"
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

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

type Generator struct {
	certbot      cert.Generator
	systemConfig cert.GeneratorSystemConfig
	logger       *log.Logger
	dateProvider date.Provider
}

type GeneratorUserConfig interface {
	IsCertbotStaging() bool
	GetUserEmail() *string
	GetDomain() *string
	GetDomainUpdateToken() *string
}

func New(systemConfig cert.GeneratorSystemConfig, dateProvider date.Provider, certbot cert.Generator) *Generator {
	return &Generator{
		certbot:      certbot,
		systemConfig: systemConfig,
		logger:       log.Default(),
		dateProvider: dateProvider,
	}
}

func (g *Generator) Start() {
	file, err := os.OpenFile("", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("unable to create certbot logger %v\n", err)
	}
	g.logger = log.New(file, "", log.LstdFlags)
}

func (g *Generator) RegenerateIfNeeded() error {

	certBytes, err := ioutil.ReadFile(g.systemConfig.SslCertificateFile())
	if err != nil {
		return err
	}

	certificate, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return err
	}

	now := g.dateProvider.Now()
	validFor := certificate.NotAfter.Sub(now)
	valid := validFor > Month
	realCert := certificate.Subject.String() != fake.Subject
	if valid && realCert {
		g.logger.Printf("not regenerating real certificate, valid for days: %d\n", int(validFor.Hours()/24))
		return nil
	}

	return g.Generate()
}

func (g *Generator) Generate() error {
	err := g.certbot.Generate()
	if err != nil {
		err = ioutil.WriteFile(Log, []byte(err.Error()), 644)
	}
	return err
}
