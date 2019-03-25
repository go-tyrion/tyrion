package http

import (
	"fmt"
	"net/http"
	"time"
)

type Options struct {
	IP           string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HandleFunc func(*Context)

type HttpServer struct {
	options *Options
	router  *Router
}

func New() *HttpServer {
	return &HttpServer{
		router: new(Router),
	}
}

// 初始化
func (this *HttpServer) Init(options *Options) {
	this.options = options
}

func (this *HttpServer) InitByConfig(confFile string) {
	this.options = this.resolveConfigToOptions(confFile)
}

// http
func (this *HttpServer) Run() error {
	addr := fmt.Sprintf("%s:%d", this.options.IP, this.options.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}

// https
func (this *HttpServer) RunAsTLS() error {

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
