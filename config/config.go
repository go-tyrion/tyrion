package config

import (
	"github.com/go-ini/ini"
	"os"
	"sync"
)

const BaseConfigPath = "config/"

var cfg *Config

func init() {
	cfg = newConfig()
}

func newConfig() *Config {
	return &Config{
		env:   os.Getenv("APP_ENV"),
		cache: make(map[string]*ini.Section),
	}
}

type Config struct {
	mux sync.Mutex

	env   string
	cache map[string]*ini.Section
}

func (c *Config) getKey(file string) *ini.Section {
	if k, ok := cfg.cache[file]; ok {
		return k
	}

	f, err := ini.Load(BaseConfigPath + file)
	if err != nil {
		panic(err)
	}

	c.mux.Lock()
	c.cache[file] = f.Section(c.env)
	c.mux.Unlock()

	return f.Section(c.env)
}

// String get value for string
func String(field, file string) string {
	return cfg.getKey(file).Key(field).String()
}

func Strings(field, file, delim string) []string {
	return cfg.getKey(file).Key(field).Strings(delim)
}

// Int get value for string
func Int(field, file string) (int, error) {
	return cfg.getKey(file).Key(field).Int()
}
