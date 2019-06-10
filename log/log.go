package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	SuffixFormatForHour = "2006010215"
	SuffixFormatForDay  = "20060102"

	DateTimeFormat = "2006-01-02 15:04:05.000"
)

type (
	LogLevel      int
	LogRotateType string
)

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	PANIC
	FATAL
	PRINT
)

const (
	RotateNone   LogRotateType = ""
	RotateHourly LogRotateType = "H"
	RotateDaily  LogRotateType = "D"
)

var levels = []string{
	"debug",
	"info",
	"warn",
	"error",
	"panic",
	"fatal",
	"print",
}

var _log *Logger

func init() {
	_log = NewLogger()
	log.Print()
}

type Logger struct {
	mu sync.Mutex

	// 日志级别, 默认 "log.Debug"
	level LogLevel

	// 日志切割方式，支持按 "D"、"H"，即 "按天"、"按小时" 进行切割
	// 默认不切割，可以通过 "log.SetRotateHourly()" 和 "log.SetRotateDaily()" 修改
	// 只有当指定以文件方式输出生效
	rotateType LogRotateType

	// 文件名，以文件方式输出
	// 文件后缀，当指定切割方式时生效
	dir, file, suffix string

	// 前缀信息
	prefix string

	// 是否显示文件和执行方法
	showCaller bool

	// 输出句柄，默认以标准方式输出
	out io.Writer

	// 输出格式，支持 "text" 和 "json" 格式输出
	formatter Formatter
}

func NewLogger() *Logger {
	l := &Logger{
		level:      DEBUG,
		rotateType: RotateNone,
		out:        os.Stdout,
	}
	l.formatter = NewTextFormatter(l)

	return l
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) ShowCaller() {
	l.showCaller = true
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) SetRotateHourly() {
	l.rotateType = RotateHourly
}

func (l *Logger) SetRotateDaily() {
	l.rotateType = RotateDaily
}

func (l *Logger) SetTextFormatter() {
	l.formatter = NewTextFormatter(l)
}

func (l *Logger) SetJsonFormatter() {
	l.formatter = NewJsonFormatter(l)
}

func (l *Logger) SetOutputDir(dir string) {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if mkErr := os.Mkdir(dir, 0644); mkErr != nil {
				panic(mkErr.Error())
			}
		} else {
			panic(err.Error())
		}
	}

	l.dir = dir
}

func (l *Logger) SetOutputByName(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.file = name
}

func (l *Logger) log(level LogLevel, dep int, v ...interface{}) {
	if l.level > level {
		return
	}

	text, err := l.formatter.Format(level, dep, fmt.Sprintln(v...))
	if err != nil {
		return
	}

	l.write(text)
}

func (l *Logger) logf(level LogLevel, dep int, f string, v ...interface{}) {
	if l.level > level {
		return
	}

	text, err := l.formatter.Format(level, dep, fmt.Sprintf(f, v...))
	if err != nil {
		return
	}

	l.write(text)
}

func (l *Logger) write(text []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.rotate()

	if _, err := l.out.Write(text); err != nil {
		fmt.Println("WErr:", err.Error())
	}
}

func (l *Logger) rotate() {
	if l.file == "" {
		return
	}

	suffix := l.genSuffix()
	if l.suffix == "" || (l.suffix != suffix && l.rotateType != RotateNone) {
		l.suffix = suffix
		l.setOutput()
	}

	return
}

func (l *Logger) Debug(v ...interface{}) {
	l.log(DEBUG, 4, v...)
}

func (l *Logger) Debugf(f string, v ...interface{}) {
	l.logf(DEBUG, 4, f, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.log(INFO, 4, v...)
}

func (l *Logger) Infof(f string, v ...interface{}) {
	l.logf(INFO, 4, f, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(WARN, 4, v...)
}

func (l *Logger) Warnf(f string, v ...interface{}) {
	l.logf(WARN, 4, f, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.log(ERROR, 4, v...)
}

func (l *Logger) Errorf(f string, v ...interface{}) {
	l.logf(ERROR, 4, f, v...)
}

func (l *Logger) Panic(v ...interface{}) {
	msg := concat(v...)
	l.log(PANIC, 4, v...)
	panic(msg)
}

func (l *Logger) Panicf(f string, v ...interface{}) {
	msg := concat(v...)
	l.logf(PANIC, 4, f, v...)
	panic(msg)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log(FATAL, 4, v...)
	os.Exit(1)
}

func (l *Logger) Fatalf(f string, v ...interface{}) {
	l.logf(FATAL, 4, f, v...)
	os.Exit(1)
}

func (l *Logger) Print(v ...interface{}) {
	l.log(PRINT, 4, v...)
}

func (l *Logger) Printf(f string, v ...interface{}) {
	l.logf(PRINT, 4, f, v...)
}

func (l *Logger) genSuffix() string {
	var suffix string

	if l.rotateType == RotateHourly {
		suffix = time.Now().Format(SuffixFormatForHour)
	} else if l.rotateType == RotateDaily {
		suffix = time.Now().Format(SuffixFormatForDay)
	}

	return suffix
}

func concat(msg ...interface{}) string {
	buf := make([]string, 0, len(msg))
	for _, m := range msg {
		buf = append(buf, fmt.Sprintf("%v", m))
	}
	return strings.Join(buf, " ")
}

func (l *Logger) setOutput() {
	var fileName string
	if l.rotateType == "" {
		fileName = filepath.Join(l.dir, l.file)
	} else {
		fileName = filepath.Join(l.dir, l.file) + "." + l.suffix
	}

	h, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err.Error())
	}

	l.out = h
}

// ------------------------------------------------------------
func SetLevel(level LogLevel) {
	_log.SetLevel(level)
}

func ShowCaller() {
	_log.ShowCaller()
}

func SetRotateHourly() {
	_log.SetRotateHourly()
}

func SetRotateDaily() {
	_log.SetRotateDaily()
}

func SetTextFormatter() {
	_log.SetTextFormatter()
}

func SetJsonFormatter() {
	_log.SetJsonFormatter()
}

func SetOutputDir(dir string) {
	_log.SetOutputDir(dir)
}

func SetOutputByName(name string) {
	_log.SetOutputByName(name)
}

func Debug(v ...interface{}) {
	_log.Debug(v...)
}

func Debugf(f string, v ...interface{}) {
	_log.Debugf(f, v...)
}

func Info(v ...interface{}) {
	_log.Info(v...)
}

func Infof(f string, v ...interface{}) {
	_log.Infof(f, v...)
}

func Warn(v ...interface{}) {
	_log.Warn(v...)
}

func Warnf(f string, v ...interface{}) {
	_log.Warnf(f, v...)
}

func Error(v ...interface{}) {
	_log.Error(v...)
}

func Errorf(f string, v ...interface{}) {
	_log.Errorf(f, v...)
}

func Panic(v ...interface{}) {
	_log.Panic(v...)
}

func Panicf(f string, v ...interface{}) {
	_log.Panicf(f, v...)
}

func Fatal(v ...interface{}) {
	_log.Fatal(v...)
}

func Fatalf(f string, v ...interface{}) {
	_log.Fatalf(f, v...)
}

func Print(v ...interface{}) {
	_log.Print(v...)
}

func Printf(f string, v ...interface{}) {
	_log.Printf(f, v...)
}
