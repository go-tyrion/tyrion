package ini

import (
	"github.com/go-ini/ini"
	"path/filepath"
)

func MapTo(file string, p interface{}) error {
	cfg, err := ini.Load(filepath.Join(BaseConfigPath, file))
	if err != nil {
		return err
	}

	return cfg.MapTo(p)
}
