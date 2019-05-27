package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	SuffixFormatForHour = "2006010215"
	SuffixFormatForDay  = "20060102"
)

type LogLevel int

const (
	LDebug LogLevel = iota
	LInfo
	LWarn
	LError
	LPanic
	LFatal
)

var levels = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"PANIC",
	"FATAL",
}

var _log *logger

func init() {
	_log = NewLogger()
}

type logger struct {
	mu sync.Mutex

	level LogLevel

	rotateHourly bool
	rotateDaily  bool

	fileName   string
	fileSuffix string
	fileHandle *os.File

	formater Formater
}

func NewLogger() *logger {
	return &logger{
		level:        LDebug,
		rotateHourly: false,
		rotateDaily:  false,
		fileHandle:   os.Stdout,
		formater:     new(TextFormater),
	}
}

func (l *logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *logger) SetRotateHourly() {
	l.rotateHourly = true
	l.rotateDaily = false
	l.fileSuffix = l.genSuffix()
}

func (l *logger) SetRotateDaily() {
	l.rotateDaily = true
	l.rotateHourly = false
	l.fileSuffix = l.genSuffix()
}

func (l *logger) SetOutputByName(name string) (err error) {
	var h *os.File
	h, err = os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}

	l.fileName = name
	l.fileHandle = h

	return
}

func (l *logger) log(level LogLevel, v ...interface{}) {
	if l.level > level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.rotate(); err != nil {
		return
	}

	msg := make([]interface{}, len(v)+2)
	msg[0] = "[" + levels[level] + "]"
	copy(msg[1:], v)
	msg[len(v)+1] = ""

	val, err := l.formater.Format(msg)
	if err != nil {
		return
	}

	l.fileHandle.Write(val)
}

func (l *logger) logf(level LogLevel, f string, v ...interface{}) {
	if l.level > level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.rotate(); err != nil {
		return
	}

	msg := "[" + levels[level] + "] " + fmt.Sprintf(f, v...)

	val, err := l.formater.Format(msg)
	if err != nil {
		return
	}

	l.fileHandle.Write(val)
}

func (l *logger) rotate() (err error) {
	if !l.rotateDaily && !l.rotateHourly {
		return
	}

	if l.fileSuffix == l.genSuffix() {
		return
	}

	newFileName := l.fileName + "." + l.fileSuffix
	err = os.Rename(l.fileName, newFileName)
	if err != nil {
		return
	}

	return l.SetOutputByName(l.fileName)
}

func (l *logger) Debug(v ...interface{}) {
	l.log(LDebug, v...)
}

func (l *logger) Debugf(f string, v ...interface{}) {
	l.logf(LDebug, f, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.log(LInfo, v...)
}

func (l *logger) Infof(f string, v ...interface{}) {
	l.logf(LInfo, f, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.log(LWarn, v...)
}

func (l *logger) Warnf(f string, v ...interface{}) {
	l.logf(LWarn, f, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.log(LError, v...)
}

func (l *logger) Errorf(f string, v ...interface{}) {
	l.logf(LError, f, v...)
}

func (l *logger) Panic(v ...interface{}) {
	l.log(LPanic, v...)
}

func (l *logger) Panicf(f string, v ...interface{}) {
	l.logf(LPanic, f, v...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.log(LFatal, v...)
	os.Exit(1)
}

func (l *logger) Fatalf(f string, v ...interface{}) {
	l.logf(LFatal, f, v...)
	os.Exit(1)
}

func (l *logger) genSuffix() (suffix string) {
	if l.rotateHourly {
		suffix = time.Now().Format(SuffixFormatForHour)
	} else if l.rotateDaily {
		suffix = time.Now().Format(SuffixFormatForDay)
	}
	return
}
