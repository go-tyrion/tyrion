package http

import "time"

const (
	DefaultMaxPostMemory = 30 << 20 // max form multipart memory size, default 32M
)

type Options struct {
	IP                  string
	Port                int
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	TLSCertFile         string
	TLSKeyFile          string
	IgnorePathLastSlash bool // 忽略路由后面的斜线
}

var (
	// 常用默认参数
	optsDefaultPort         = 8001
	optsDefaultReadTimeout  = time.Duration(30) * time.Second
	optsDefaultWriteTimeout = time.Duration(30) * time.Second
)

func (o *Options) DefaultOpts() *Options {
	return o.ResetOpts(nil)
}

func (o *Options) ResetOpts(opts *Options) *Options {
	if opts.Port <= 0 {
		opts.Port = optsDefaultPort
	}

	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = optsDefaultReadTimeout
	}

	if opts.WriteTimeout == 0 {
		opts.WriteTimeout = optsDefaultWriteTimeout
	}

	return opts
}

func (o *Options) ResolveOptsByConfigFile(name string) *Options {
	return nil
}
