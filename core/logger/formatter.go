package log

type Formatter interface {
	Format(level LogLevel, v interface{}) (b []byte, err error)
}
