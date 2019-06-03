package log

import (
	log2 "github.com/ngaut/log"
	"log"
	"os"
	"testing"
)

func TestNormal(t *testing.T) {
	logger := log.New(os.Stdout, "prefix ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Println("message", "message2")
	// log.Panicln("message")
	//log.Fatalln("message")
}

func TestQiNiuLog(t *testing.T) {
	log2.Info("info1", "info2")
	_log.SetJsonFormatter()
	_log.Debug("debug1", "debug2")
	_log.Info("info1", "info2")
}

func TestLogger_Info(t *testing.T) {
	// _log.SetOutputByName("app.log")
	// _log.SetRotateHourly()

	_log.SetPrefix("[Tyrion]")

	_log.SetOutputByName("demo.log")
	_log.SetJsonFormatter()

	for i := 0; i < 10000; i++ {
		_log.Info("this is info message, i:", i, "message2", "message3")
		// time.Sleep(time.Second * 1)
	}
}
