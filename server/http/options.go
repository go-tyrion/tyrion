package http

import "time"

type Options struct {
	Addr                string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	TLSCertFile         string
	TLSKeyFile          string
	MaxPostMemory       int64
	IgnorePathLastSlash bool // 忽略路由后面的斜线
}

func (o *Options) DefaultOpts() *Options {
	return o.ResetOpts(nil)
}

func (o *Options) ResetOpts(opts *Options) *Options {
	if opts.Addr == "" {
		opts.Addr = ":8080"
	}

	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = time.Duration(30) * time.Second
	}

	if opts.WriteTimeout == 0 {
		opts.WriteTimeout = time.Duration(30) * time.Second
	}

	if opts.MaxPostMemory <= 0 {
		opts.MaxPostMemory = 30 << 20
	}

	return opts
}

func (opt *Options) GetMaxPostMemory() int64 {
	return opt.MaxPostMemory
}

func (o *Options) ResolveOptsByConfigFile(name string) *Options {
	return nil
}
