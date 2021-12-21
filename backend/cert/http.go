package cert

import (
	"os"
)

type HttpProviderSyncloud struct{}

const Path = "/var/snap/platform/current/certbot/www/.well-known/acme-challenge/"

func NewHttpProviderSyncloud() *HttpProviderSyncloud {
	return &HttpProviderSyncloud{}
}

func (d *HttpProviderSyncloud) Present(_, token, keyAuth string) error {
	path := Path + token
	err := os.WriteFile(path, []byte(keyAuth), 0644)
	return err
}

func (d *HttpProviderSyncloud) CleanUp(_, token, keyAuth string) error {
	return os.Remove(Path + token)
}
