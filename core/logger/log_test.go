package log

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	logger := log.New(os.Stdout, "prefix ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Println("message")
	TestLog()
	// log.Panicln("message")
	//log.Fatalln("message")
}

func TestLogger_Info(t *testing.T) {
	_log.SetOutputByName("app.log")
	_log.SetRotateHourly()

	for i := 0; i < 10000; i++ {
		_log.Infof("this is info message, i: %d", i)
		time.Sleep(time.Second * 1)
	}
}
