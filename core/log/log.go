package log

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// 功能规划
// 1. 常用日志输出，区分Info, Error等日志级别，支持到文件（不同文件）
// 2. 支持把日志输出到其他的 server 里面，如filebeat, logstash,syslog,rsyslog,Fluentd 等
// 3. 支持 json 输出
// 4. 支持添加公共的字段
// 5. 支持日志切分
// 6. traceID

type LogLevel int

const (
	Ldebug LogLevel = iota
	Linfo
	Lwarn
	Lerror
	Lpanic
	Lfatal
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
	_log = new(logger)
	_log.init()
}

type RotationType string

const (
	RotationHourly RotationType = "H"
	RotationDaily  RotationType = "D"
)

const (
	LogFormatJSON = iota
	LogFormatText
)

type File struct {
	l      *log.Logger
	file   string
	handle *os.File
}

func (f *File) GetFile() string {
	return f.file
}

func (f *File) GetHandle() *os.File {
	return f.handle
}

type logger struct {
	mu sync.Mutex

	// distinguishFile bool             // 是否区分文件
	// file            *File            // 不区分文件时保存文件句柄
	files map[int]*File // 区分文件时保存文件句柄

	level    LogLevel
	rotation RotationType
}

func New() *logger {
	return &logger{
		level: Ldebug,
		files: make(map[int]*File),
	}
}

func (l *logger) init() {
	for index, level := range levels {
		f := new(File)
		f.l = log.New(os.Stderr, "", log.LstdFlags)
		f.file = level + ".log"

		l.files[index] = new(File)
	}
}

func (l *logger) output(level LogLevel, v []interface{}) {
	if level < l.level { // 低于这个级别的日志不输出
		return
	}

	vl := make([]interface{}, len(v)+2)
	vl[0] = "[" + levels[level] + "]"
	copy(vl[1:0], v)
	vl[len(v)+1] = ""

	l.files[int(level)].l.Output(4, fmt.Sprintln(vl...))
}

func (l *logger) rotate() {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch l.rotation {
	case RotationHourly:

	case RotationDaily:

	}
}

// 是否分文件保存
func (l *logger) DistinguishFile(dist bool) {
	// l.distinguishFile = dist
}

// todo file handle
func (l *logger) SetLevelFile(level LogLevel, file string) {
	if len(levels) <= int(level) {
		h, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}

		// l.SetOutput(f)

		f := new(File)
		f.file = file
		f.handle = h

		l.mu.Lock()
		l.files[int(level)] = f
		l.mu.Unlock()
	}
}

func (l *logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *logger) SetRotation(r RotationType) {
	l.rotation = r
}

func (l *LogLevel) Debug(v ...interface{}) {

}

func (l *logger) Info(v ...interface{}) {

}

func (l *logger) Warn(v ...interface{}) {

}

func (l *logger) Error(v ...interface{}) {

}

func (l *logger) Panic(v ...interface{}) {

}

func (l *logger) Fatal(v ...interface{}) {

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 函数
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DistFile(files map[string]string) {

}

func SetRotation(r RotationType) {
	_log.SetRotation(r)
}

func SetLevel(level LogLevel) {
	_log.SetLevel(level)
}

func Info(v ...interface{}) {
	_log.output(v)
}
