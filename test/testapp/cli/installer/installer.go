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
	AuthUrl         string
	AuthLocalSocket string
	AppUrl          string
	ClientSecret    string
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

	clientSecret, err := client.RegisterOIDCClient(App, []string{"/oidc/callback"}, false, "client_secret_basic")
	if err != nil {
		return err
	}

	err = config.Generate(
		path.Join(AppDir, "config"),
		path.Join(DataDir, "config"),
		Variables{
			AuthUrl:         authUrl,
			AuthLocalSocket: client.GetAuthLocalSocket(),
			AppUrl:          appUrl,
			ClientSecret:    clientSecret,
		},
	)
	if err != nil {
		return err
	}

	err = createSnapUserLinks()
	if err != nil {
		return err
	}

	err = linux.Chown(DataDir, App)
	if err != nil {
		return err
	}
	return linux.Chown(CommonDir, App)
}

func createSnapUserLinks() error {
	revision := os.Getenv("SNAP_REVISION")
	snapUserDir := path.Join(CommonDir, App, "snap", App)
	err := os.MkdirAll(path.Join(snapUserDir, revision), 0755)
	if err != nil {
		return err
	}
	err = replaceSymlink(revision, path.Join(snapUserDir, "current"))
	if err != nil {
		return err
	}
	return replaceSymlink("missing", path.Join(CommonDir, "dangling"))
}

func replaceSymlink(target, link string) error {
	err := os.Remove(link)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Symlink(target, link)
}
