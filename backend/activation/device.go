package activation

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/cert"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/nginx"
	"log"
	"strings"
)

type Device struct {
	config               DevicePlatformUserConfig
	certificateGenerator *cert.Generator
	auth                 *auth.Service
	nginx                *nginx.Nginx
	trigger              *event.Trigger
	cookies              Cookies
}

type Cookies interface {
	Reset()
}

type DevicePlatformUserConfig interface {
	GetRedirectDomain() string
	SetActivated()
	SetWebSecretKey(key string)
	SetIpv4Public(enabled bool)
	SetPublicPort(port *int)
	SetPublicIp(publicIp *string)
}

type DeviceActivation interface {
	ActivateDevice(username string, password string, name string, email string) error
}

func NewDevice(
	config DevicePlatformUserConfig,
	auth *auth.Service,
	nginx *nginx.Nginx,
	trigger *event.Trigger,
	cookies Cookies,
) *Device {
	return &Device{
		config:  config,
		auth:    auth,
		nginx:   nginx,
		trigger: trigger,
		cookies: cookies,
	}
}

func (d *Device) ActivateDevice(username string, password string, name string, email string) error {
	err := d.resetAccess()
	if err != nil {
		return err
	}

	d.config.SetWebSecretKey(uuid.New().String())
	d.cookies.Reset()

	err = d.auth.Reset(name, username, password, email)
	if err != nil {
		return err
	}

	err = d.nginx.InitConfig()
	if err != nil {
		return err
	}
	err = d.nginx.ReloadPublic()
	if err != nil {
		return err
	}

	d.config.SetActivated()

	log.Println("activation completed")

	return nil
}

func (d *Device) resetAccess() error {
	log.Println("reset access")
	d.config.SetIpv4Public(false)
	d.config.SetPublicIp(nil)
	d.config.SetPublicPort(nil)
	return d.trigger.RunAccessChangeEvent()
}

func ParseUsername(username string, domain string) (string, string) {
	if strings.Contains(username, "@") {
		return strings.Split(username, "@")[0], username
	}
	email := fmt.Sprintf("%s@%s", username, domain)
	return username, email
}
