package core

const DefaultAppConfigFileName = "app.ini"

type App struct {
	configFile string
}

func LoadAppConfig() {
	app := new(App)
	app.configFile = ""
}
