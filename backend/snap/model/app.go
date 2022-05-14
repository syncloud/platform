package model

import (
	"fmt"
)

type App struct {
	Name string `json:"name"`
	Snap string `json:"snap"`
}

func (app *App) RunCommand() string {
	return fmt.Sprintf("%v.%v", app.Snap, app.Name)
}
