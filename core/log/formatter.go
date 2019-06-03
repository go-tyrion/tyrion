package log

type Formatter interface {
	Format(level LogLevel, v string) (b []byte, err error)
}
