package config

import (
	"github.com/go-ini/ini"
	"lib/config/proto"
	"lib/helper"
	"os"
	"path/filepath"
	"sync"
)

const BaseConfigPath = "config"

var _instance *Config

func init() {
	_instance = newConfig()
	_instance.init()
}

func newConfig() *Config {
	return &Config{
		section: ini.DefaultSection,
		cache:   make(map[string]*ini.File),
	}
}

type Config struct {
	mux sync.Mutex

	section string
	cache   map[string]*ini.File
}

func (c *Config) init() {
	app := new(proto.AppConfig)
	if err := c.Resolve("app", app); err != nil {
		panic(err)
	}

	if app.Env != "" {
		c.section = app.Env
	} else {
		c.section = "prod"
	}

	c.section = app.Env

	_ = os.Setenv("env", c.section)
	_ = os.Setenv("debug", helper.Bool2String(app.Debug))
}

// 将配置与数据结构映射
func (c *Config) Resolve(file string, p interface{}) error {
	cfg, err := ini.Load(getFullPath(file))
	if err != nil {
		return err
	}

	cfg.NameMapper = ini.TitleUnderscore

	defaultSection := cfg.Section(ini.DefaultSection)
	// env keys
	envKeys := cfg.Section(c.section).KeyStrings()
	for _, key := range envKeys {
		value := cfg.Section(c.section).Key(key).Value()
		if defaultSection.HasKey(key) {
			defaultSection.Key(key).SetValue(value)
		} else {
			_, _ = defaultSection.NewKey(key, value)
		}
	}

	return defaultSection.MapTo(p)
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

func (c *Config) GetKey(file string, field string) *ini.Key {
	f := c.getFile(file)
	s := f.Section(c.section)

	if s.HasKey(field) {
		return s.Key(field)
	}

	return f.Section(ini.DefaultSection).Key(field)
}

func getFullPath(file string) string {
	wishedExt := ".ini"

	ext := filepath.Ext(file)
	if ext != wishedExt {
		if ext != "" {
			panic("unsupported config file, ext:" + ext)
		} else {
			file += wishedExt
		}
	}

	return filepath.Join(BaseConfigPath, file)
}
