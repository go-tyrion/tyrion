package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

type Formatter interface {
	Format(level LogLevel, dep int, v string) (b []byte, err error)
}

// TextFormatter 文本
type TextFormatter struct {
	logger *logger
}

func NewTextFormatter(logger *logger) Formatter {
	return &TextFormatter{
		logger: logger,
	}
}

func (f *TextFormatter) Format(level LogLevel, dep int, v string) (b []byte, err error) {
	now := time.Now()

	var text bytes.Buffer
	text.WriteString(f.logger.prefix)
	text.WriteString(now.Format(DateTimeFormat) + " ")

	if f.logger.showFile {
		fileAndLine := getCaller(dep)
		text.WriteString(fileAndLine + " ")
	}

	text.WriteString("[" + levels[level] + "]: ")
	text.WriteString(v)

	if len(v) == 0 || v[len(v)-1] != '\n' {
		text.WriteString("\n")
	}

	return text.Bytes(), nil
}

// JsonFormatter Json 格式
type JsonFormatter struct {
	logger *logger
}

func NewJsonFormatter(l *logger) *JsonFormatter {
	return &JsonFormatter{
		logger: l,
	}
}

func (f *JsonFormatter) Format(level LogLevel, dep int, v string) (b []byte, err error) {
	things := make(map[string]interface{})

	msgLength := len(v)
	if msgLength > 0 && v[msgLength-1] == '\n' {
		v = v[0 : msgLength-1]
	}

	things["time"] = time.Now().Format(DateTimeFormat)
	things["level"] = levels[level]
	things["message"] = v

	if f.logger.showFile {
		things["file"] = getCaller(dep)
	}

	thingsBuffer := &bytes.Buffer{}

	encoder := json.NewEncoder(thingsBuffer)
	if err := encoder.Encode(things); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}

	return thingsBuffer.Bytes(), nil
}

func getCaller(dep int) string {
	file, line, ok, _ := "--", 0, true, "--"
	_, file, line, ok = runtime.Caller(dep)
	if ok {
		// file = filepath.Base(file)
		file += ":" + strconv.Itoa(line)
	}

	return file
}
