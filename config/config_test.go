package config

import (
	"os"
	"testing"
)

func TestString(t *testing.T) {
	e := os.Setenv("APP_ENV", "prod")
	if e != nil {
		t.Fatal(e)
	}
	s := String("app.name", "config.ini")
	t.Log(s)
}
