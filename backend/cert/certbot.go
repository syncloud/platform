package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"os"
)

type Certbot struct {
	redirect     RedirectCertbot
	userConfig   UserConfig
	systemConfig GeneratorSystemConfig
}

type UserConfig interface {
	IsCertbotStaging() bool
	GetUserEmail() *string
	GetDomain() *string
	GetDomainUpdateToken() *string
}

type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}
func (u User) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

type CertbotGenerator interface {
	Generate() error
}

func NewCertbot(redirect RedirectCertbot, userConfig UserConfig, systemConfig GeneratorSystemConfig) *Certbot {
	return &Certbot{
		redirect:     redirect,
		userConfig:   userConfig,
		systemConfig: systemConfig,
	}
}

func (g *Certbot) Generate() error {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	email := g.userConfig.GetUserEmail()
	if email == nil {
		return fmt.Errorf("email is not set")
	}
	myUser := User{
		Email: *email,
		key:   privateKey,
	}

	certbotConfig := lego.NewConfig(&myUser)
	if g.userConfig.IsCertbotStaging() {
		certbotConfig.CADirURL = lego.LEDirectoryStaging
	}

	client, err := lego.NewClient(certbotConfig)
	if err != nil {
		return err
	}

	if g.userConfig.GetDomain() {
		token := g.userConfig.GetDomainUpdateToken()
		if token == nil {
			return fmt.Errorf("token is not set")
		}
		err = client.Challenge.SetDNS01Provider(NewDNSProviderSyncloud(*token, g.redirect))
		if err != nil {
			return err
		}
	} else {
		err = client.Challenge.SetHTTP01Provider(NewHttpProviderSyncloud())
		if err != nil {
			return err
		}
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return err
	}
	myUser.Registration = reg

	domain := g.userConfig.GetDomain()
	if domain == nil {
		return fmt.Errorf("domain is not set")
	}
	request := certificate.ObtainRequest{
		Domains: []string{
			*domain,
			fmt.Sprintf("*.%s", *domain),
		},
		Bundle: true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	certificateFile := g.systemConfig.SslCertificateFile()
	err = os.WriteFile(certificateFile, certificates.Certificate, 0644)
	if err != nil {
		return err
	}

	keyFile := g.systemConfig.SslKeyFile()
	err = os.WriteFile(keyFile, certificates.PrivateKey, 0644)
	if err != nil {
		return err
	}

	return nil
}
