package installer

import (
	"github.com/syncloud/golib/config"
	"github.com/syncloud/golib/linux"
	"github.com/syncloud/golib/platform"
	"os"
	"path"
)

const (
	App       = "testapp"
	AppDir    = "/snap/testapp/current"
	DataDir   = "/var/snap/testapp/current"
	CommonDir = "/var/snap/testapp/common"
)

type Variables struct {
	AuthUrl      string
	AppUrl       string
	ClientSecret string
}

func Install() error {
	return UpdateConfigs()
}

func Configure() error {
	return UpdateConfigs()
}

func AccessChange() error {
	return UpdateConfigs()
}

func StorageChange() error {
	client := platform.New()
	_, err := client.InitStorage(App, App)
	return err
}

func UpdateConfigs() error {
	err := linux.CreateUser(App)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(DataDir, "config"), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(DataDir, "nginx"), 0755)
	if err != nil {
		return err
	}

	client := platform.New()
	authUrl, err := client.GetAppUrl("auth")
	if err != nil {
		return err
	}

	appUrl, err := client.GetAppUrl(App)
	if err != nil {
		return err
	}

	clientSecret, err := client.RegisterOIDCClient(App, "/oidc/callback", false, "client_secret_basic")
	if err != nil {
		return err
	}

	err = config.Generate(
		path.Join(AppDir, "config"),
		path.Join(DataDir, "config"),
		Variables{
			AuthUrl:      authUrl,
			AppUrl:       appUrl,
			ClientSecret: clientSecret,
		},
	)
	if err != nil {
		return err
	}

	err = linux.Chown(DataDir, App)
	if err != nil {
		return err
	}
	return linux.Chown(CommonDir, App)
}
