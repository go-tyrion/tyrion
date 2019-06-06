package config

import (
	"testing"
)

func TestString(t *testing.T) {
	s := String("app.name", "config.ini")
	t.Log(s)
}
