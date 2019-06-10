package http

import (
	"lib/config"
	"lib/config/proto"
	"time"
)

func newOptions(file string) *Options {
	opts := new(Options)
	opts.Init(file)

	return opts
}

type Options struct {
	proto.HttpConfig

	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IgnorePathLastSlash bool // 忽略路由后面的斜线
}

func (opts *Options) Init(file string) {
	err := config.Resolve(file, &opts.HttpConfig)
	if err != nil {
		panic(err)
	}
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

	return opts
}

func (opt *Options) GetMaxPostMemory() int64 {
	return 0
}

func (o *Options) ResolveOptsByConfigFile(name string) *Options {
	return nil
}
