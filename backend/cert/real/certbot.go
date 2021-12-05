package real

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/syncloud/platform/cert"
	"os"
)

type Certbot struct {
	redirect     RedirectCertbot
	userConfig   GeneratorUserConfig
	systemConfig cert.GeneratorSystemConfig
}

func NewCertbot(redirect RedirectCertbot, userConfig GeneratorUserConfig, systemConfig cert.GeneratorSystemConfig) *Certbot {
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
	myUser := MyUser{
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

	token := g.userConfig.GetDomainUpdateToken()
	if token == nil {
		return fmt.Errorf("token is not set")
	}
	err = client.Challenge.SetDNS01Provider(NewDNSProviderSyncloud(*token, g.redirect))
	if err != nil {
		return err
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
