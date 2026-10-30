package caddy

import "github.com/syncloud/platform/config"

type UserConfigMock struct {
	deviceDomain string
	port         *int
}

func (u *UserConfigMock) Url(app string) string {
	return config.ConstructUrl(u.port, app+"."+u.deviceDomain)
}

func (u *UserConfigMock) GetDeviceDomain() string {
	return u.deviceDomain
}
