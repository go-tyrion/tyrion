package proto

type AppConfig struct {
	Name  string `ini:"name"`
	Env   string `ini:"env"`
	Debug bool   `ini:"debug"`
}
