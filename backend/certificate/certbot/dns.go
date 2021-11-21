package certbot

import "github.com/go-acme/lego/v4/challenge/dns01"

type DNSProviderSyncloud struct {
	apiAuthToken string
	redirect     RedirectCertbot
}

type RedirectCertbot interface {
	CertbotPresent(domain, txtValue string) error
	CertbotCleanUp(domain string) error
}

func NewDNSProviderSyncloud(apiAuthToken string, redirect RedirectCertbot) *DNSProviderSyncloud {
	return &DNSProviderSyncloud{
		apiAuthToken: apiAuthToken,
		redirect:     redirect,
	}
}

func (d *DNSProviderSyncloud) Present(domain, _, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	// make API request to set a TXT record on fqdn with value and TTL
	return d.redirect.CertbotPresent(fqdn, value)
}

func (d *DNSProviderSyncloud) CleanUp(domain, token, keyAuth string) error {
	// clean up any state you created in Present, like removing the TXT record
	return d.redirect.CertbotCleanUp(domain)
}
