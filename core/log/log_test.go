package log

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	logger := log.New(os.Stdout, "prefix ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Println("message", "message2")
	// log.Panicln("message")
	//log.Fatalln("message")
}

func TestQiNiuLog(t *testing.T) {
	_logger := NewLogger()
	_logger.ShowCaller(true)
	_logger.Info("lhh0")
}

func TestLogger_Info(t *testing.T) {
	_log.SetPrefix("[Tyrion]")
	_log.SetOutputByName("demo.log")
	_log.SetJsonFormatter()
	_log.SetRotateHourly()
	_log.ShowCaller(true)

	for i := 0; i < 10000; i++ {
		_log.Info("this is info message, i:", i, "message2", "message3")
		time.Sleep(time.Second * 1)
	}
}
