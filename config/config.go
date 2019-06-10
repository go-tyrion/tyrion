package config

import (
	"lib/config/ini"
	"os"
)

func init() {
	ini.SetEnv(os.Getenv("env"))
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

func Int64(field, file string) int64 {
	val, _ := ini.GetKey(file, field).Int64()
	return val
}

// Bool
func Bool(field, file string) bool {
	val, _ := ini.GetKey(file, field).Bool()
	return val
}
