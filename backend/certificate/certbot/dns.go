package certbot

import "github.com/go-acme/lego/v4/challenge/dns01"

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
	// make API request to set a TXT record on fqdn with value and TTL
	return d.redirect.CertbotPresent(d.token, fqdn, value)
}

func (d *DNSProviderSyncloud) CleanUp(domain, _, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	// clean up any state you created in Present, like removing the TXT record
	return d.redirect.CertbotCleanUp(d.token, fqdn, value)
}
