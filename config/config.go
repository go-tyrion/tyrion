package config

import (
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"sync"
)

const BaseConfigPath = "config"

var _cfg *Config

func getInstance() *Config {
	if _cfg == nil {
		_cfg = newConfig()
	}

	return _cfg
}

func newConfig() *Config {
	return &Config{
		env:   os.Getenv("APP_ENV"),
		cache: make(map[string]*ini.File),
	}
}

type Config struct {
	mux sync.Mutex

	env   string
	cache map[string]*ini.File
}

func (c *Config) getFile(file string) *ini.File {
	if k, ok := getInstance().cache[file]; ok {
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

func (c *Config) getKey(file string, field string) *ini.Key {
	f := c.getFile(file)
	s := f.Section(c.env)

	if s.HasKey(field) {
		return s.Key(field)
	}

	return f.Section(ini.DefaultSection).Key(field)
}

// String get value for string
func String(field, file string) string {
	return getInstance().getKey(file, field).String()
}

func Strings(field, file, delim string) []string {
	return getInstance().getKey(file, field).Strings(delim)
}

// Int get value for string
func Int(field, file string) int {
	val, _ := getInstance().getKey(file, field).Int()
	return val
}

// Bool
func Bool(field, file string) bool {
	val, _ := getInstance().getKey(file, field).Bool()
	return val
}
