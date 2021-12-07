package cert

import (
	"github.com/go-acme/lego/v4/challenge/dns01"
	"time"
)

type DNSProviderSyncloud struct {
	token    string
	redirect RedirectCertbot
	values   []string
}

type RedirectCertbot interface {
	CertbotPresent(token, fqdn string, value ...string) error
	CertbotCleanUp(token, fqdn string) error
}

func NewDNSProviderSyncloud(token string, redirect RedirectCertbot) *DNSProviderSyncloud {
	return &DNSProviderSyncloud{
		token:    token,
		redirect: redirect,
	}
}

func (d *DNSProviderSyncloud) Present(domain, _, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	d.values = append(d.values, value)
	return d.redirect.CertbotPresent(d.token, fqdn, d.values...)
}

func (d *DNSProviderSyncloud) CleanUp(domain, _, keyAuth string) error {
	d.values = make([]string, 0)
	fqdn, _ := dns01.GetRecord(domain, keyAuth)
	return d.redirect.CertbotCleanUp(d.token, fqdn)
}

func (d *DNSProviderSyncloud) Timeout() (timeout, interval time.Duration) {
	return 5 * time.Minute, 60 * time.Second
}