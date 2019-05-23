package log

import (
	"log"
	"os"
	"strings"
	"sync"
)

// 功能规划
// 1. 常用日志输出，区分Info, Error等日志级别，支持到文件（不同文件）
// 2. 支持把日志输出到其他的 server 里面，如filebeat, logstash,syslog,rsyslog,Fluentd 等
// 3. 支持 json 输出
// 4. 支持添加公共的字段
// 5. 支持日志切分
// 6. traceID

const (
	Ldebug = iota
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

	distinguishFile bool             // 是否区分文件
	file            *File            // 不区分文件时保存文件句柄
	files           map[string]*File // 区分文件时保存文件句柄

	l        *log.Logger
	level    int
	rotation RotationType
}

func New() *logger {
	return &logger{}
}

func (l *logger) init() {
	l.files = make(map[string]*File)
	for _, level := range levels {
		l.files[strings.ToLower(level)] = new(File)
	}
}

func (l *logger) output(v []interface{}) {

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
	l.distinguishFile = dist
}

// todo file handle
func (l *logger) SetLevelFile(level int, file string) {
	if len(levels) <= level {
		h, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}

		// l.SetOutput(f)

		f := new(File)
		f.file = file
		f.handle = h

		l.mu.Lock()
		l.files[levels[level]] = f
		l.mu.Unlock()
	}
}

func (l *logger) SetLevel(level int) {
	l.level = level
}

func (l *logger) SetRotation(r RotationType) {
	l.rotation = r
}

func (l *logger) Info(v ...interface{}) {

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 函数
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DistFile(files map[string]string) {

}

func SetRotation(r RotationType) {
	_log.SetRotation(r)
}

func Info(v ...interface{}) {
	_log.output(v)
}
