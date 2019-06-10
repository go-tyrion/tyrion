package proto

type HttpConfig struct {
	ServiceName    string `ini:"service_name"`
	Addr           string `ini:"addr"`
	AccessLog      bool   `ini:"access_log"`
	ReadTimeoutMs  int64  `ini:"read_timeout_ms"`
	WriteTimeoutMs int64  `ini:"write_timeout_ms"`
	MaxPostMemory  string `ini:"max_post_memory"`
}
