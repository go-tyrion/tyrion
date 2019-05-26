package log

import "testing"

func TestLogger_Info(t *testing.T) {
	_log.SetOutputByFileName("app.log")

	for i := 0; i < 10000; i++ {
		_log.Infof("this is info message, i: %d", i)
	}
}
