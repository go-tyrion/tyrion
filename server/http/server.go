package http

import (
	"fmt"
	"net/http"
)

type Options struct {
	IP   string
	Port int
	// ReadTimeout  time.Duration
	// WriteTimeout time.Duration
}

type HandleFunc func(*Context)

type HttpServer struct {
	opts   *Options
	router *Router
}

func New() *HttpServer {
	return &HttpServer{
		router: new(Router),
	}
}

// 初始化
func (this *HttpServer) Init(opts *Options) {
	this.opts = opts
}

// 通过配置文件初始化
func (this *HttpServer) InitByConfig(confFile string) {
	this.opts = this.resolveConfigToOptions(confFile)
}

// http
func (this *HttpServer) Run() error {
	addr := fmt.Sprintf("%s:%d", this.opts.IP, this.opts.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}

// https
func (this *HttpServer) RunTLS() error {

	return nil
}

// ------------
// 类型 bisinessServer模式
func (this *HttpServer) AddLogic(path string, logic string) {

}

// 传统方式
func (this *HttpServer) Get(path string, h Context) {

}

func (this *HttpServer) Post(path string, h Context) {

}

// DI

// USE

// ------------
// 私有方法
func (this *HttpServer) resolveConfigToOptions(confFile string) *Options {

	return nil
}
