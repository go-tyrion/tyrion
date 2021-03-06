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
	logger *Logger
}

func NewTextFormatter(logger *Logger) Formatter {
	return &TextFormatter{
		logger: logger,
	}
}

func (f *TextFormatter) Format(level LogLevel, dep int, v string) (b []byte, err error) {
	now := time.Now()

	var text bytes.Buffer
	text.WriteString(f.logger.prefix)
	text.WriteString(now.Format(DateTimeFormat) + " ")

	if f.logger.showCaller {
		fileAndLine := getCaller(dep)
		text.WriteString(fileAndLine + " ")
	}

	if level != PRINT {
		text.WriteString("[" + levels[level] + "]: ")
	}

	text.WriteString(v)

	if len(v) == 0 || v[len(v)-1] != '\n' {
		text.WriteString("\n")
	}

	return text.Bytes(), nil
}

// JsonFormatter Json 格式
type JsonFormatter struct {
	logger *Logger
}

func NewJsonFormatter(l *Logger) *JsonFormatter {
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
	things["message"] = v

	if level != PRINT {
		things["level"] = levels[level]
	}

	if f.logger.showCaller {
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
