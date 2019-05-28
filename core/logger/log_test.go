package log

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	logger := log.New(os.Stdout, "prefix ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Println("message")
	// log.Panicln("message")
	//log.Fatalln("message")
}

func TestLogrus(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Debug("debug1", "debug2")
	logrus.Info("hello")
	logrus.Warn("warn1", "warn2")
}

func TestLogger_Info(t *testing.T) {
	_log.SetOutputByName("app.log")
	_log.SetRotateHourly()

	for i := 0; i < 10000; i++ {
		_log.Infof("this is info message, i: %d", i)
		time.Sleep(time.Second * 1)
	}
}
