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

const (
	DefaultMaxPostMemory = 30 << 20 // max form multipart memory size, default 32M
)

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

type HandleFunc func(ctx *Context)

type HttpServer struct {
	debug         bool
	opts          *Options
	router        *Router
	server        *http.Server
	logger        *log.Logger
	pool          sync.Pool
	maxPostMemory int64
}

func New() *HttpServer {
	app := new(HttpServer)
	app.router = NewRouter(app)
	app.server = new(http.Server)
	app.logger = log.New(os.Stdout, "[Tyrion] ", log.LstdFlags)
	app.pool = sync.Pool{
		New: func() interface{} {
			return NewContext(nil, nil, app)
		},
	}
	app.maxPostMemory = DefaultMaxPostMemory
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

func (s *HttpServer) GetMaxPostMemory() int64 {
	return s.maxPostMemory
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
	s.add(http.MethodGet, pattern, h)
}

func (s *HttpServer) Post(path string, h Context) {

}

func (s *HttpServer) add(method string, pattern string, handles []HandleFunc) {
	wrapHandles := make([]HandleFunc, 0, len(handles))
	for _, h := range handles {
		wrapHandles = append(wrapHandles, WrapHandlerFunc(h))
	}
	s.router.Add(method, pattern, wrapHandles)
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := s.pool.Get().(*Context)
	ctx.reset(w, r)

	handles := s.router.Get(r.URL.Path)
	if handles == nil {
		ctx.handles = append(ctx.handles, defaultNotFound())
	} else {
		ctx.handles = append(ctx.handles, handles...)
	}

	ctx.Next()

	s.pool.Put(ctx)
}

func (s *HttpServer) Stop() error {
	return s.server.Close()
}

// DI

// USE

// ------------
// 私有方法
// default 404
func defaultNotFound() HandleFunc {
	return func(ctx *Context) {
		ctx.String(404, "not found!")
		ctx.Next()
	}
}

// WrapHandleFunc wrap for context handler chain
func WrapHandlerFunc(h HandleFunc) HandleFunc {
	return func(ctx *Context) {
		h(ctx)
		ctx.Next()
	}
}

func (s *HttpServer) resolveConfigToOptions(confFile string) *Options {

	return nil
}
