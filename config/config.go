package config

import (
	"lib/config/ini"
	"lib/core"
)

func init() {
	ini.SetEnv(core.Env())
}

// String get value for string
func String(field, file string) string {
	return ini.GetKey(file, field).String()
}

func Strings(field, file, delim string) []string {
	return ini.GetKey(file, field).Strings(delim)
}

// Int get value for string
func Int(field, file string) int {
	val, _ := ini.GetKey(file, field).Int()
	return val
}

// Bool
func Bool(field, file string) bool {
	val, _ := ini.GetKey(file, field).Bool()
	return val
}
