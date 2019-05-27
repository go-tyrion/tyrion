package log

type Formater interface {
	Format(v interface{}) (b []byte, err error)
}
