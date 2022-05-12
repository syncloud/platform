package snap

import "fmt"

type Snap struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Channel string `json:"channel"`
	Version string `json:"version"`
	Type    string `json:"type"`
	Apps    []App  `json:"apps"`
}

type App struct {
	Name string `json:"name"`
	Snap string `json:"snap"`
}

func (s *Snap) ToSyncloudApp(url string) SyncloudAppVersions {
	return SyncloudAppVersions{
		CurrentVersion: &s.Version,
		App: SyncloudApp{
			Id:   s.Name,
			Name: s.Summary,
			Url:  url,
			Icon: fmt.Sprintf("/rest/app_image?channel=%s&app=%s", s.Channel, s.Name),
		},
	}
}

func (s *Snap) FindApp(app string) (bool, *App) {
	for _, snapApp := range s.Apps {
		if snapApp.Name == app {
			return true, &snapApp
		}
	}
	return false, nil
}

func (app *App) RunCommand() string {
	return fmt.Sprintf("%v.%v", app.Snap, app.Name)
}
