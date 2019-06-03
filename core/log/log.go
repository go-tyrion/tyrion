package log

import (
	"fmt"
	"io"
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
	LDebug LogLevel = iota
	LInfo
	LWarn
	LError
	LPanic
	LFatal
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
}

var _log *logger

func init() {
	_log = NewLogger()
}

type logger struct {
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

func NewLogger() *logger {
	l := &logger{
		level:      LDebug,
		rotateType: RotateNone,
		out:        os.Stdout,
	}
	l.formatter = NewTextFormatter(l)

	return l
}

func (l *logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *logger) ShowCaller(show bool) {
	l.showCaller = show
}

func (l *logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *logger) SetRotateHourly() {
	l.rotateType = RotateHourly
	l.suffix = l.genSuffix()
}

func (l *logger) SetRotateDaily() {
	l.rotateType = RotateDaily
	l.suffix = l.genSuffix()
}

func (l *logger) SetTextFormatter() {
	l.formatter = NewTextFormatter(l)
}

func (l *logger) SetJsonFormatter() {
	l.formatter = NewJsonFormatter(l)
}

func (l *logger) SetOutputDir(dir string) {
	l.dir = dir
}

func (l *logger) SetOutputByName(name string) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.file = name
	l.suffix = l.genSuffix()

	var h *os.File
	h, err = os.OpenFile(l.buildFile(), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err.Error())
	}

	l.out = h

	return
}

func (l *logger) buildFile() string {
	return filepath.Join(l.dir, l.file)
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

	text := fmt.Sprintln(v...)

	val, err := l.formatter.Format(level, text)
	if err != nil {
		return
	}

	l.out.Write(val)
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

	msg := fmt.Sprintf(f, v...)

	val, err := l.formatter.Format(level, msg)
	if err != nil {
		return
	}

	l.out.Write(val)
}

func (l *logger) rotate() (err error) {
	if l.rotateType == "" {
		return
	}

	if l.suffix == l.genSuffix() {
		return
	}

	newFileName := l.file + "." + l.suffix
	err = os.Rename(l.file, newFileName)
	if err != nil {
		return
	}

	return l.SetOutputByName(l.file)
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
	msg := concat(v...)
	l.log(LPanic, v...)
	panic(msg)
}

func (l *logger) Panicf(f string, v ...interface{}) {
	msg := concat(v...)
	l.logf(LPanic, f, v...)
	panic(msg)
}

func (l *logger) Fatal(v ...interface{}) {
	l.log(LFatal, v...)
	os.Exit(1)
}

func (l *logger) Fatalf(f string, v ...interface{}) {
	l.logf(LFatal, f, v...)
	os.Exit(1)
}

func (l *logger) genSuffix() string {
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

// ------------------------------------------------------------
func SetLevel(level LogLevel) {
	_log.SetLevel(level)
}

func ShowCaller(show bool) {
	_log.ShowCaller(show)
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
