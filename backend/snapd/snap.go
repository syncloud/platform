package snapd

import "fmt"

type Snap struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Channel string `json:"channel"`
	Version string `json:"version"`
	Apps    []App  `json:"apps"`
}

type App struct {
	Name string `json:"name"`
	Snap string `json:"snap"`
}

func (snap *Snap) FindApp(app string) (bool, *App) {
	for _, snapApp := range snap.Apps {
		if snapApp.Name == app {
			return true, &snapApp
		}
	}
	return false, nil
}

func (app *App) RunCommand() string {
	return fmt.Sprintf("snap run %v.%v", app.Snap, app.Name)
}
