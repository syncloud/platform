package config

import (
	"github.com/bigkevmcd/go-configparser"
)

const WebCertificatePort = 80
const WebAccessPort = 443
const WebProtocol = "https"

type SystemConfig struct {
	parser *configparser.ConfigParser
}

const File = "/snap/platform/current/config/platform.cfg"

func NewSystemConfig(file string) (*SystemConfig, error) {
	parser, err := configparser.NewConfigParserFromFile(file)
	if err != nil {
		return nil, err
	}

	config := &SystemConfig{
		parser: parser,
	}
	return config, nil
}

func (c *SystemConfig) DataDir() (*string, error) {
	return c.get("data_dir")
}

func (c *SystemConfig) AppDir() (*string, error) {
	return c.get("app_dir")
}

func (c *SystemConfig) ConfigDir() (*string, error) {
	return c.get("config_dir")
}

func (c *SystemConfig) get(key string) (*string, error) {
	value, err := c.parser.Get("platform", key)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
