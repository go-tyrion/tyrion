package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Options struct {
	IP   string
	Port int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	TLSCertFile string
	TLSKeyFile  string
}

type HandleFunc func(*Context)

type HttpServer struct {
	opts   *Options
	router *Router
	server *http.Server
	addr   string
}

func New() *HttpServer {
	return &HttpServer{
		router: new(Router),
		server: new(http.Server),
	}
}

// 初始化
func (s *HttpServer) Init(opts *Options) {
	s.opts = opts
	s.addr = fmt.Sprintf("%s:%d", s.opts.IP, s.opts.Port)
	s.server.ReadTimeout = s.opts.ReadTimeout
	s.server.WriteTimeout = s.opts.WriteTimeout
}

// 通过配置文件初始化
func (s *HttpServer) InitByConfig(confFile string) {
	s.opts = s.resolveConfigToOptions(confFile)
}

// http
func (s *HttpServer) Run() error {
	s.server.Handler = s
	return s.server.ListenAndServe()
}

// https
func (s *HttpServer) RunTLS() error {
	if s.opts.TLSCertFile == "" || s.opts.TLSKeyFile == "" {
		return errors.New("invalid tls config")
	}

	s.server.Handler = s
	return s.server.ListenAndServeTLS(s.opts.TLSCertFile, s.opts.TLSKeyFile)
}

// ------------
// 类型 bisinessServer模式
func (s *HttpServer) AddLogic(path string, logic string) {

}

// 传统方式
func (s *HttpServer) Get(path string, h Context) {
	s.router.Add(GET, path, nil)
}

func (s *HttpServer) Post(path string, h Context) {

}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

// DI

// USE

// ------------
// 私有方法
func (s *HttpServer) resolveConfigToOptions(confFile string) *Options {

	return nil
}
