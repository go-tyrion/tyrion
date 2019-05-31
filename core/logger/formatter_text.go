package log

import (
	"bytes"
	"fmt"
	"time"
)

type TextFormatter struct {
	logger *logger
}

func NewTextFormatter(logger *logger) Formatter {
	return &TextFormatter{
		logger: logger,
	}
}

func (f *TextFormatter) Format(level LogLevel, v interface{}) (b []byte, err error) {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("input must be string")
	}

	now := time.Now()

	var text bytes.Buffer
	text.WriteString(f.logger.prefix)
	text.WriteString(now.Format(DateTimeFormat) + " ")
	text.WriteString("[" + levels[level] + "]: ")
	text.WriteString(s)

	if len(s) == 0 || s[len(s)-1] != '\n' {
		text.WriteString("\n")
	}

	return text.Bytes(), nil
}
