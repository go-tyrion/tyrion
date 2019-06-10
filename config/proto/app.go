package proto

type AppConfig struct {
	Name   string `ini:"name"`
	Env    string `ini:"env"`
	Debug  bool   `ini:"debug"`
	LogDir string `ini:"log_dir"`
}
