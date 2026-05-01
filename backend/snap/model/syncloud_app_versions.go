package model

type SyncloudAppVersions struct {
	App              SyncloudApp `json:"app"`
	CurrentVersion   *string     `json:"current_version"`
	InstalledVersion *string     `json:"installed_version"`
	LocalInstall     bool        `json:"local_install"`
}
