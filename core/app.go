package core

import (
	"lib/config/ini"
	"lib/config/proto"
	"os"
)

const DefaultAppConfigFileName = "app.ini"

type App struct {
	proto.AppConfig
}

func (app *App) Init() {
	err := ini.MapTo(DefaultAppConfigFileName, &app.AppConfig)
	if err != nil {
		panic(app)
	}

	_ = os.Setenv("env", app.Env)
}
