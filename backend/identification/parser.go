package identification

import (
	"github.com/bigkevmcd/go-configparser"
	"log"
	"net"
)

const DefaultIdFile = "/etc/syncloud/id.cfg"

type Id struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	MacAddress string `json:"mac_address"`
}

type Parser struct {
	filename string
}

type IdParser interface {
	Id() (*Id, error)
}

func New() *Parser {
	return &Parser{filename: DefaultIdFile}
}

func (p *Parser) get(key string, def string) string {
	config, err := configparser.NewConfigParserFromFile(p.filename)
	if err != nil {
		log.Printf("cannot load id config: %s, %s", p.filename, err)
		config = configparser.New()
	}

	option, err := config.HasOption("id", key)
	if err != nil {
		log.Printf("identification key (%s) error: %s", key, err)
		return def
	}
	if option {
		option, err := config.Get("id", key)
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
