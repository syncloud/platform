package snap

type SyncloudApp struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Required    string `json:"required"`
	Ui          string `json:"ui"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type Apps struct {
	Apps []SyncloudApp `json:"apps"`
}

type SyncloudAppVersions struct {
	App              SyncloudApp `json:"app"`
	CurrentVersion   *string     `json:"current_version"`
	InstalledVersion *string     `json:"installed_version"`
}
