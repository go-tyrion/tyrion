package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
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

var (
	// 常用默认参数
	defaultHttpServerOpts = func() *Options {
		return &Options{
			Port:         8001,
			ReadTimeout:  time.Duration(30) * time.Second,
			WriteTimeout: time.Duration(30) * time.Second,
		}
	}
)

type HandleFunc func(*Context)

type HttpServer struct {
	debug  bool
	opts   *Options
	router *Router
	server *http.Server
	logger *log.Logger
	pool   sync.Pool
}

func New() *HttpServer {
	app := new(HttpServer)
	app.router = NewRouter(app)
	app.server = new(http.Server)
	app.logger = log.New(os.Stdout, "[Tyrion] ", log.LstdFlags)
	app.pool = sync.Pool{
		New: func() interface{} {
			return NewContext(nil, nil)
		},
	}
	return app
}

// 使用 Default 默认配置
func Default() *HttpServer {
	app := New()
	app.Init(defaultHttpServerOpts())
	return app
}

// 通过 Init 方法初始化
func (s *HttpServer) Init(opts *Options) {
	s.opts = opts
	s.server.Addr = fmt.Sprintf("%s:%d", s.opts.IP, s.opts.Port)
	s.server.ReadTimeout = s.opts.ReadTimeout
	s.server.WriteTimeout = s.opts.WriteTimeout
}

// 通过配置文件初始化
func (s *HttpServer) InitByConfig(confFile string) {
	s.Init(s.resolveConfigToOptions(confFile))
}

func (s *HttpServer) Log() *log.Logger {
	return s.logger
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
func (s *HttpServer) Get(pattern string, h ...HandleFunc) {
	s.router.Add(http.MethodGet, pattern, h)
}

func (s *HttpServer) Post(path string, h Context) {

}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := s.pool.Get().(*Context)
	ctx.Reset(w, r)

	handles := s.router.Get(r.URL.Path)
	if handles == nil {
		return
	} else {
		ctx.handles = append(ctx.handles, handles...)
	}

	ctx.Run()

	s.pool.Put(ctx)
}

// DI

// USE

// ------------
// 私有方法
func (s *HttpServer) resolveConfigToOptions(confFile string) *Options {

	return nil
}
