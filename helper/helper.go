package helper

import (
	"unicode"
)

// UCFirst 首字母大写
func UCFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// LCFirst 首字母小写
func LCFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func Bool2String(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}
