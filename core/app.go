package core

import (
	"lib/config/ini"
)

const DefaultAppConfigFileName = "app.ini"

var _app *app

func init() {
	_app = new(app)
	_app.loadConfig()
}

type app struct {
	name  string
	env   string
	debug bool
}

func (a *app) loadConfig() {
	a.name = ini.GetKey(DefaultAppConfigFileName, "app.name").String()
	a.env = ini.GetKey(DefaultAppConfigFileName, "app.env").String()
	a.debug, _ = ini.GetKey(DefaultAppConfigFileName, "app.debug").Bool()
}

func (a *app) getName() string {
	return a.name
}

func (a *app) getEnv() string {
	return a.env
}

func (a *app) getDebug() bool {
	return a.debug
}

// ------------------------------------------
func Name() string {
	return _app.getName()
}

func Env() string {
	return _app.getEnv()
}

func Debug() bool {
	return _app.getDebug()
}
