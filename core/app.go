package core

import (
	"lib/config"
	"lib/config/proto"
)

type App struct {
	proto.AppConfig
}

func (app *App) Init() {
	err := config.Resolve(config.DefaultAppConfigFile, &app.AppConfig)
	if err != nil {
		panic(app)
	}
}
