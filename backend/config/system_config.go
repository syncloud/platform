package config

import (
	"github.com/bigkevmcd/go-configparser"
	"log"
)

const WebAccessPort = 443
const WebProtocol = "https"

type SystemConfig struct {
	file   string
	parser *configparser.ConfigParser
}

var DefaultSystemConfig string

func init() {
	DefaultSystemConfig = "/snap/platform/current/config/platform.cfg"
}

func NewSystemConfig(file string) *SystemConfig {
	return &SystemConfig{
		file: file,
	}
}

func (c *SystemConfig) Load() {
	parser, err := configparser.NewConfigParserFromFile(c.file)
	if err != nil {
		log.Fatalln(err)
	}
	c.parser = parser
}

func (c *SystemConfig) DataDir() string {
	return c.get("data_dir")
}

func (c *SystemConfig) CommonDir() string {
	return c.get("common_dir")
}

func (c *SystemConfig) AppDir() string {
	return c.get("app_dir")
}

func (c *SystemConfig) ConfigDir() string {
	return c.get("config_dir")
}

func (c *SystemConfig) SslCertificateFile() string {
	return c.get("ssl_certificate_file")
}

func (c *SystemConfig) SslKeyFile() string {
	return c.get("ssl_key_file")
}

func (c *SystemConfig) SslCaCertificateFile() string {
	return c.get("ssl_ca_certificate_file")
}

func (c *SystemConfig) SslCaKeyFile() string {
	return c.get("ssl_ca_key_file")
}

func (c *SystemConfig) Channel() string {
	return c.get("channel")
}

func (c *SystemConfig) DiskLink() string {
	return c.get("disk_link")
}

func (c *SystemConfig) ExternalDiskDir() string {
	return c.get("external_disk_dir")
}

func (c *SystemConfig) InternalDiskDir() string {
	return c.get("internal_disk_dir")
}

func (c *SystemConfig) get(key string) string {
	value, err := c.parser.GetInterpolated("platform", key)
	if err != nil {
		log.Fatal(err)
	}
	return value
}
