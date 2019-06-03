package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type JsonFormatter struct {
	logger *logger
}

func NewJsonFormatter(l *logger) *JsonFormatter {
	return &JsonFormatter{
		logger: l,
	}
}

func (f *JsonFormatter) Format(level LogLevel, v string) (b []byte, err error) {
	things := make(map[string]interface{})

	msgLength := len(v)
	if msgLength > 0 && v[msgLength-1] == '\n' {
		v = v[0 : msgLength-1]
	}

	things["time"] = time.Now().Format(DateTimeFormat)
	things["level"] = levels[level]
	things["message"] = v

	thingsBuffer := &bytes.Buffer{}

	encoder := json.NewEncoder(thingsBuffer)
	if err := encoder.Encode(things); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}

	return thingsBuffer.Bytes(), nil
}
