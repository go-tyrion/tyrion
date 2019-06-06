package ini

import (
	"github.com/go-ini/ini"
	"path/filepath"
	"sync"
)

const BaseConfigPath = "config"

var _cfg *Config

func init() {
	_cfg = newConfig()
}

func newConfig() *Config {
	return &Config{
		env:   ini.DefaultSection,
		cache: make(map[string]*ini.File),
	}
}

type Config struct {
	mux sync.Mutex

	env   string
	cache map[string]*ini.File
}

func (c *Config) setEnv(env string) {
	c.env = env
}

func (c *Config) getFile(file string) *ini.File {
	if k, ok := c.cache[file]; ok {
		return k
	}

	f, err := ini.Load(filepath.Join(BaseConfigPath, file))
	if err != nil {
		panic(err)
	}

	c.mux.Lock()
	c.cache[file] = f
	c.mux.Unlock()

	return f
}

func SetEnv(env string) {
	_cfg.setEnv(env)
}

func GetKey(file string, field string) *ini.Key {
	f := _cfg.getFile(file)
	s := f.Section(_cfg.env)

	if s.HasKey(field) {
		return s.Key(field)
	}

	return f.Section(ini.DefaultSection).Key(field)
}
