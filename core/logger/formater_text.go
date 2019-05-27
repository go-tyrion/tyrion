package log

import (
	"bytes"
	"fmt"
)

type TextFormater struct {
}

func (f *TextFormater) Format(v interface{}) (b []byte, err error) {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("input must be string")
	}

	a := bytes.NewBufferString(s)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		a.WriteString("\n")
	}

	return a.Bytes(), nil
}
