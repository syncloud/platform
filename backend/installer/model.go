package installer

type AppInstaller interface {
	Upgrade(version string) error
}
