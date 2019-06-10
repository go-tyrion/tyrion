package proto

type HttpConfig struct {
	ServiceName     string
	Addr            string
	AccessLog       bool
	AccessLogRotate string
	ReadTimeoutMs   int64
	WriteTimeoutMs  int64
	MaxPostMemory   string
	HttpsCertFile   string
	HttpsKeyFile    string
}
