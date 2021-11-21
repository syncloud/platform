package activation

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/certificate/selfsigned"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/nginx"
	"log"
	"strings"
)

type Device struct {
	config               DevicePlatformUserConfig
	certificateGenerator *selfsigned.Generator
	auth                 *auth.Service
	nginx                *nginx.Nginx
	trigger              *event.Trigger
}

type DevicePlatformUserConfig interface {
	GetRedirectDomain() string
	SetActivated()
	SetWebSecretKey(key string)
	SetExternalAccess(enabled bool)
	SetUpnp(enabled bool)
	SetManualCertificatePort(manualCertificatePort int)
	SetManualAccessPort(manualAccessPort int)
	DeletePublicIp()
}

type DeviceActivation interface {
	ActivateDevice(username string, password string, name string, email string) error
}

func NewDevice(
	config DevicePlatformUserConfig,
	auth *auth.Service,
	nginx *nginx.Nginx,
	trigger *event.Trigger,
) *Device {
	return &Device{
		config:  config,
		auth:    auth,
		nginx:   nginx,
		trigger: trigger,
	}
}

func (d *Device) ActivateDevice(username string, password string, name string, email string) error {
	err := d.resetAccess()
	if err != nil {
		return err
	}

	d.config.SetWebSecretKey(uuid.New().String())

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
	d.config.SetUpnp(false)
	d.config.SetExternalAccess(false)
	d.config.DeletePublicIp()
	d.config.SetManualCertificatePort(0)
	d.config.SetManualAccessPort(0)
	return d.trigger.RunAccessChangeEvent()
}

func ParseUsername(username string, domain string) (string, string) {
	if strings.Contains(username, "@") {
		return strings.Split(username, "@")[0], username
	}
	email := fmt.Sprintf("%s@%s", username, domain)
	return username, email
}
