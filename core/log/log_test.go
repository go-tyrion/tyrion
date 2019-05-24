package log

import "testing"

func TestSetLevel(t *testing.T) {
	_log.SetOutputType(OutputFile)
	_log.SetOutputFile("app.log")
	_log.Info("info message.")
}
