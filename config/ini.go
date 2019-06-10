package config

import (
	"github.com/go-ini/ini"
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
	env := c.getFile("app").Section(c.section).Key("env").String()
	if env != "" {
		c.section = env
	} else {
		c.section = "prod"
	}

	c.section = env
	_ = os.Setenv("env", c.section)
}

// 将配置与数据结构映射
func (c *Config) Resolve(file string, p interface{}) error {
	cfg, err := ini.Load(getFullPath(file))
	if err != nil {
		return err
	}

	cfg.NameMapper = ini.TitleUnderscore

	newSection, _ := cfg.NewSection(file)
	// default keys
	defaultKeys := cfg.Section(ini.DefaultSection).KeyStrings()
	for _, key := range defaultKeys {
		_, _ = newSection.NewKey(key, cfg.Section(ini.DefaultSection).Key(key).Value())
	}

	// env keys
	envKeys := cfg.Section(c.section).KeyStrings()
	for _, key := range envKeys {
		value := cfg.Section(c.section).Key(key).Value()
		if newSection.HasKey(key) {
			newSection.Key(key).SetValue(value)
		} else {
			_, _ = newSection.NewKey(key, value)
		}
	}

	return newSection.MapTo(p)
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
	ext := filepath.Ext(file)
	if ext == "" {
		file += ".ini"
	} else if ext != "ini" {
		panic("unsupported config file")
	}

	return filepath.Join(BaseConfigPath, file)
}
