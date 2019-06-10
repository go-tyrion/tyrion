package config

// String get value for string
func String(field, file string) string {
	return _instance.GetKey(file, field).String()
}

func Strings(field, file, delim string) []string {
	return _instance.GetKey(file, field).Strings(delim)
}

// Int get value for string
func Int(field, file string) int {
	val, _ := _instance.GetKey(file, field).Int()
	return val
}

func Int64(field, file string) int64 {
	val, _ := _instance.GetKey(file, field).Int64()
	return val
}

// Bool
func Bool(field, file string) bool {
	val, _ := _instance.GetKey(file, field).Bool()
	return val
}

func Resolve(file string, p interface{}) error {
	return _instance.Resolve(file, p)
}
