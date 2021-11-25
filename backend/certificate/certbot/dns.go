package certbot

import (
	"github.com/go-acme/lego/v4/challenge/dns01"
	"time"
)

type DNSProviderSyncloud struct {
	token    string
	redirect RedirectCertbot
}

type RedirectCertbot interface {
	CertbotPresent(token, fqdn, value string) error
	CertbotCleanUp(token, fqdn, value string) error
}

func NewDNSProviderSyncloud(token string, redirect RedirectCertbot) *DNSProviderSyncloud {
	return &DNSProviderSyncloud{
		token:    token,
		redirect: redirect,
	}
}

func (d *DNSProviderSyncloud) Present(domain, _, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	return d.redirect.CertbotPresent(d.token, fqdn, value)
}

func (d *DNSProviderSyncloud) CleanUp(domain, _, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	return d.redirect.CertbotCleanUp(d.token, fqdn, value)
}

func (d *DNSProviderSyncloud) Timeout() (timeout, interval time.Duration) {
	return 5 * time.Minute, 1 * time.Minute
}

