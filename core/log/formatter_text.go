package log

import (
	"bytes"
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

func (f *TextFormatter) Format(level LogLevel, v string) (b []byte, err error) {
	now := time.Now()

	var text bytes.Buffer
	text.WriteString(f.logger.prefix)
	text.WriteString(now.Format(DateTimeFormat) + " ")
	text.WriteString("[" + levels[level] + "]: ")
	text.WriteString(v)

	if len(v) == 0 || v[len(v)-1] != '\n' {
		text.WriteString("\n")
	}

	return text.Bytes(), nil
}
