package log

type Formatter interface {
	Format(v interface{}) (b []byte, err error)
}
