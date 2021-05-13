package identification

import (
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"log"
	"net"
)

type Id struct {
	name       string
	title      string
	macAddress string
}

type Parser struct {
	config *configparser.ConfigParser
}

func New(filename string) (*Parser, error) {

	config, err := configparser.NewConfigParserFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot load id config: %s, %s", filename, err)
	}

	return &Parser{config: config}, nil
}

func (p *Parser) get(key string, def string) string {
	option, err := p.config.HasOption("id", key)
	if err != nil {
		log.Printf("identification key (%s) error: %s", key, err)
		return def
	}
	if option {
		option, err := p.config.Get("id", key)
		if err != nil {
			log.Printf("identification key (%s) error: %s", key, err)
			return def
		}
		return option
	}
	return def
}

func (p *Parser) name() string {
	return p.get("name", "unknown")
}

func (p *Parser) title() string {
	return p.get("title", "Unknown")
}

func (p *Parser) Id() (*Id, error) {
	mac, err := GetMac()
	if err != nil {
		return nil, err
	}
	return &Id{p.name(), p.title(), mac}, nil
}

func GetMac() (string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, ifa := range ifas {
		addr := ifa.HardwareAddr.String()
		if len(ifa.HardwareAddr) >= 6 && ifa.Name != "" {
			return addr, nil
		}
	}
	return "", nil

}
