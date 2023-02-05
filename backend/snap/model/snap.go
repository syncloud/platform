package model

import (
	"fmt"
)

type Snap struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Channel string `json:"channel"`
	Version string `json:"version"`
	Type    string `json:"type"`
	Apps    []App  `json:"apps"`
}

func (s *Snap) ToStoreApp(url string) SyncloudAppVersions {
	app := s.toSyncloudApp(url)
	app.CurrentVersion = &s.Version
	return app
}

func (s *Snap) ToInstalledApp(url string) SyncloudAppVersions {
	app := s.toSyncloudApp(url)
	app.InstalledVersion = &s.Version
	return app
}

func (s *Snap) toSyncloudApp(url string) SyncloudAppVersions {
	return SyncloudAppVersions{
		App: SyncloudApp{
			Id:   s.Name,
			Name: s.Summary,
			Url:  url,
			Icon: fmt.Sprintf("/rest/proxy/image?channel=%s&app=%s", s.Channel, s.Name),
		},
	}
}

func (s *Snap) IsApp() bool {
	return s.Type == "app"
}

func (s *Snap) FindCommand(name string) *App {
	for _, snapApp := range s.Apps {
		if snapApp.Name == name {
			return &snapApp
		}
	}
	return nil
}
