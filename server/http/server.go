package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
)

type HandleFunc func(c *Context)

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
	s := new(HttpServer)
	s.router = NewRouter(s)
	s.server = new(http.Server)
	s.opts = new(Options)
	s.logger = log.New(os.Stdout, "[Tyrion] ", log.LstdFlags)
	s.pool = sync.Pool{
		New: func() interface{} {
			return NewContext(nil, nil, s)
		},
	}
	s.maxPostMemory = DefaultMaxPostMemory
	return s
}

// 使用 Default 默认配置
func Default() *HttpServer {
	s := New()
	s.Init(s.opts.DefaultOpts())
	return s
}

// 通过 Init 方法初始化
func (s *HttpServer) Init(opts *Options) {
	s.opts = s.opts.ResetOpts(opts)
	s.server.Addr = fmt.Sprintf("%s:%d", s.opts.IP, s.opts.Port)
	s.server.ReadTimeout = s.opts.ReadTimeout
	s.server.WriteTimeout = s.opts.WriteTimeout
}

// 通过配置文件初始化
func (s *HttpServer) InitByConfig(confFile string) {
	s.Init(s.opts.ResolveOptsByConfigFile(confFile))
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
func (s *HttpServer) AddLogic(prefix string, logic Logic) {
	t := reflect.TypeOf(logic)
	v := reflect.ValueOf(logic)
	for i := 0; i < t.NumMethod(); i++ {
		funcName := t.Method(i).Name
		if "Init" == funcName {
			continue
		}

		// todo 1 要判断首字母的情况
		// todo 2 需要判断方法签名情况

		s.add(http.MethodGet, prefix+"/"+funcName, []HandleFunc{s.wrapLogic(v.Method(i))})
	}
}

func (s *HttpServer) wrapLogic(v reflect.Value) HandleFunc {
	return func(c *Context) {
		v.Call([]reflect.Value{reflect.ValueOf(c)})
		c.Next()
	}
}

// 传统方式
func (s *HttpServer) Get(pattern string, h ...HandleFunc) {
	s.add(http.MethodGet, pattern, h)
}

func (s *HttpServer) Head(pattern string, h ...HandleFunc) {
	s.add(http.MethodHead, pattern, h)
}

func (s *HttpServer) Post(pattern string, h ...HandleFunc) {
	s.add(http.MethodPost, pattern, h)
}

func (s *HttpServer) Put(pattern string, h ...HandleFunc) {
	s.add(http.MethodPut, pattern, h)
}

func (s *HttpServer) Patch(pattern string, h ...HandleFunc) {
	s.add(http.MethodPatch, pattern, h)
}

func (s *HttpServer) Delete(pattern string, h ...HandleFunc) {
	s.add(http.MethodDelete, pattern, h)
}

func (s *HttpServer) Connect(pattern string, h ...HandleFunc) {
	s.add(http.MethodConnect, pattern, h)
}

func (s *HttpServer) Options(pattern string, h ...HandleFunc) {
	s.add(http.MethodOptions, pattern, h)
}

func (s *HttpServer) Trace(pattern string, h ...HandleFunc) {
	s.add(http.MethodTrace, pattern, h)
}

func (s *HttpServer) add(method string, pattern string, handles []HandleFunc) {
	wrapHandles := make([]HandleFunc, 0, len(handles))
	for _, h := range handles {
		wrapHandles = append(wrapHandles, WrapHandlerFunc(h))
	}
	s.router.Register(method, pattern, wrapHandles)
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := s.pool.Get().(*Context)
	c.reset(w, r)

	defer s.pool.Put(c)

	if _, ok := HttpMethods[r.Method]; !ok {
		c.handles = append(c.handles, catchHandles(405))
		return
	}

	handles := s.router.Get(r.Method, r.URL.Path)
	if handles == nil {
		c.handles = append(c.handles, catchHandles(404))
	} else {
		c.handles = handles
	}

	c.Next()
}

func (s *HttpServer) Stop() error {
	return s.server.Close()
}

// DI

// USE

// ------------
// 私有方法
// default 404

func catchHandles(code int) HandleFunc {
	return func(c *Context) {
		c.String(code, HttpStatus[code])
		c.Next()
	}
}

// WrapHandleFunc wrap for context handler chain
func WrapHandlerFunc(h HandleFunc) HandleFunc {
	return func(c *Context) {
		h(c)
		c.Next()
	}
}
