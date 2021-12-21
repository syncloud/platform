package cert

import (
	"fmt"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"go.uber.org/zap"
	"time"
)

type SyncloudDNS struct {
	token         string
	redirect      RedirectCertbot
	values        []string
	certbotLogger *zap.Logger
}

type RedirectCertbot interface {
	CertbotPresent(token, fqdn string, value ...string) error
	CertbotCleanUp(token, fqdn string) error
}

func NewSyncloudDNS(token string, redirect RedirectCertbot, certbotLogger *zap.Logger) *SyncloudDNS {
	return &SyncloudDNS{
		token:         token,
		redirect:      redirect,
		certbotLogger: certbotLogger,
	}
}

func (d *SyncloudDNS) Present(domain, _, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	d.values = append(d.values, value)
	err := d.redirect.CertbotPresent(d.token, fqdn, d.values...)
	if err != nil {
		d.certbotLogger.Error(fmt.Sprintf("dns present error: %s", err.Error()))
	}
	return err
}

func (d *SyncloudDNS) CleanUp(domain, _, keyAuth string) error {
	d.values = make([]string, 0)
	fqdn, _ := dns01.GetRecord(domain, keyAuth)
	err := d.redirect.CertbotCleanUp(d.token, fqdn)
	if err != nil {
		d.certbotLogger.Error(fmt.Sprintf("dns cleanup error: %s", err.Error()))
	}
	return err
}

func (d *SyncloudDNS) Timeout() (timeout, interval time.Duration) {
	return 5 * time.Minute, 60 * time.Second
}
