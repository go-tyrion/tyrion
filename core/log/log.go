package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
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

// 输出类型
type OutputType int

const (
	OutputConsole OutputType = iota
	OutputFile
	OutputLevelFile
)

var _log *logger

func init() {
	_log = new(logger)
	// _log.init()
}

type RotateType string

const (
	RotateHourly RotateType = "H"
	RotateDaily  RotateType = "D"
)

const (
	LogFormatJSON = iota
	LogFormatText
)

type Writer struct {
	sync.Mutex
	_logger *log.Logger

	file   string
	suffix string
	handle *os.File
}

func (w *Writer) GetLogger() *log.Logger {
	return w._logger
}

func (w *Writer) SetFile(file string) {
	w.file = file
}

func (w *Writer) GetFile() string {
	return w.file
}

func (w *Writer) GetHandle() *os.File {
	return w.handle
}

func (w *Writer) Output() {

}

func (w *Writer) Rotate(rotateType RotateType) error {
	w.Lock()
	defer w.Unlock()

	var suffix string
	if rotateType == RotateHourly {
		suffix = time.Now().Format("2006010215")
	} else if rotateType == RotateDaily {
		suffix = time.Now().Format("20060102")
	}

	fmt.Println("w.suffix:", w.suffix, " w:", suffix)

	if w.suffix != suffix {
		fmt.Println("do rotate")
		if err := w.doRotate(); err != nil {
			return err
		}

		w.suffix = suffix
	}

	return nil
}

func (w *Writer) doRotate() error {
	newFileName := w.file + "." + w.suffix
	if err := os.Rename(w.file, newFileName); err != nil {
		return err
	}

	if err := w.setOutput(); err != nil {
		return err
	}

	return w.handle.Close()
}

func (w *Writer) setOutput() error {
	f, err := os.OpenFile(w.file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	w._logger = log.New(f, w._logger.Prefix(), w._logger.Flags())
	w.handle = f

	return nil
}

type logger struct {
	mu sync.Mutex

	outputType OutputType // 输出类型
	needRotate bool       // 是否需要切割
	rotateType RotateType // 分割方式

	writer  *Writer
	writers map[int]*Writer // 区分文件时保存文件句柄

	level LogLevel
}

func New() *logger {
	l := &logger{
		level:      Ldebug,
		writers:    make(map[int]*Writer),
		outputType: OutputConsole,
		needRotate: false,
	}

	return l
}

func (l *logger) output(level LogLevel, v ...interface{}) {
	if level < l.level { // 低于这个级别的日志不输出
		return
	}

	vl := make([]interface{}, len(v)+2)
	vl[0] = "[" + levels[level] + "]"
	copy(vl[1:], v)
	vl[len(v)+1] = ""

	if l.outputType == OutputFile {
		fmt.Println("here")
		l.writer.Rotate(l.rotateType)
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

		f := new(Writer)
		f.file = file
		f.handle = h

		l.mu.Lock()
		l.writers[int(level)] = f
		l.mu.Unlock()
	}
}

// 设置日志级别
func (l *logger) SetLevel(level LogLevel) {
	l.level = level
}

// 是否需要切割
// 默认按天切割
func (l *logger) NeedRotate(need bool) {
	l.needRotate = need
	if l.rotateType == "" {
		l.rotateType = RotateDaily
	}
}

// 当 needRotate 为 true 时生效
func (l *logger) SetRotateType(r RotateType) {
	l.rotateType = r
}

// 通过名字设置分割类型
// 只接受 H 和 D 两个参数
func (l *logger) SetRotateTypeByName(name string) {
	rotateType := RotateType(name)
	if rotateType != RotateHourly && rotateType != RotateDaily {
		log.Fatal("unsupported rotate type")
	}

	l.rotateType = rotateType
}

// console / file / levelFile
// default console
func (l *logger) SetOutputType(outType OutputType) {
	l.outputType = outType
}

// 总的输出文件
func (l *logger) SetOutputFile(file string) {
	w := new(Writer)
	w.file = file

	l.writer = w
}

func (l *logger) SetOutputFileForLevel(level LogLevel, file string) {

}

func (l *logger) Debug(v ...interface{}) {

}

func (l *logger) Info(v ...interface{}) {
	l.output(Linfo, v)
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

func SetRotation(r RotateType) {
	_log.SetRotateType(r)
}

func SetLevel(level LogLevel) {
	_log.SetLevel(level)
}

func Info(v ...interface{}) {

}
