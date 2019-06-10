package core

import (
	"lib/config"
	"lib/config/proto"
)

const DefaultAppConfigFileName = "app.ini"

type App struct {
	proto.AppConfig
}

func (app *App) Init() {
	err := config.Resolve(DefaultAppConfigFileName, &app.AppConfig)
	if err != nil {
		panic(app)
	}
}
