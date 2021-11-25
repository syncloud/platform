package certbot

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
	redirect RedirectCertbot
	config   GeneratorConfig
}

type GeneratorConfig interface {
	IsCertbotStaging() bool
}

func New(redirect RedirectCertbot, config GeneratorConfig) *Generator {
	return &Generator{
		redirect: redirect,
		config:   config,
	}
}

func (g *Generator) Generate(email string, domain string, token string) error {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	myUser := MyUser{
		Email: email,
		key:   privateKey,
	}

	certbotConfig := lego.NewConfig(&myUser)
	if g.config.IsCertbotStaging() {
		certbotConfig.CADirURL = lego.LEDirectoryStaging
	}

	client, err := lego.NewClient(certbotConfig)
	if err != nil {
		return err
	}

	err = client.Challenge.SetDNS01Provider(NewDNSProviderSyncloud(token, g.redirect))
	if err != nil {
		return err
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return err
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{fmt.Sprintf("*.%s", domain)},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	err = os.WriteFile("/var/snap/platform/current/syncloud.crt", certificates.Certificate, 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile("/var/snap/platform/current/syncloud.key", certificates.PrivateKey, 0644)
	if err != nil {
		return err
	}

	return nil
}
