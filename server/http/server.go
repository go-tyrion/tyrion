package http

import (
	"lib/log"
	"net/http"
	"reflect"
	"sync"
)

type HandleFunc func(c *Context)

type HttpServer struct {
	opts   *Options
	router *Router
	server *http.Server
	logger *log.Logger
	pool   sync.Pool
}

func NewHttpServer() *HttpServer {
	server := &HttpServer{
		logger: log.NewLogger(),
		server: new(http.Server),
		opts:   new(Options),
	}
	server.router = newRouter(server)
	server.pool.New = func() interface{} {
		return newContext(server)
	}

	return server
}

// 通过 Init 方法初始化
func (s *HttpServer) Init(opts *Options) {
	s.opts = s.opts.ResetOpts(opts)
	s.server.Addr = opts.Addr
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

// Run http server
func (server *HttpServer) Run() error {
	server.setServerOpts()
	return server.server.ListenAndServe()
}

// Run https server
func (server *HttpServer) RunTLS() error {
	if server.opts.TLSCertFile == "" || server.opts.TLSKeyFile == "" {
		panic("invalid tls config")
	}

	server.setServerOpts()
	return server.server.ListenAndServeTLS(server.opts.TLSCertFile, server.opts.TLSKeyFile)
}

func (server *HttpServer) setServerOpts() {
	server.server.WriteTimeout = server.opts.WriteTimeout
	server.server.ReadTimeout = server.opts.ReadTimeout
	server.server.Handler = server
	// server.server.TLSConfig = tls.NewConfig("", "")
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

		s.Any(prefix+"/"+funcName, []HandleFunc{s.wrapLogic(v.Method(i))}...)
	}
}

func (s *HttpServer) wrapLogic(v reflect.Value) HandleFunc {
	return func(c *Context) {
		v.Call([]reflect.Value{reflect.ValueOf(c)})
		c.Next()
	}
}

// 传统方式
func (s *HttpServer) Any(pattern string, h ...HandleFunc) {
	s.add(http.MethodGet, pattern, h)
	s.add(http.MethodPost, pattern, h)
	s.add(http.MethodPut, pattern, h)
	s.add(http.MethodDelete, pattern, h)
}

func (s *HttpServer) Get(pattern string, h ...HandleFunc) {
	s.add(http.MethodGet, pattern, h)
}

func (s *HttpServer) Post(pattern string, h ...HandleFunc) {
	s.add(http.MethodPost, pattern, h)
}

func (s *HttpServer) Put(pattern string, h ...HandleFunc) {
	s.add(http.MethodPut, pattern, h)
}

func (s *HttpServer) Delete(pattern string, h ...HandleFunc) {
	s.add(http.MethodDelete, pattern, h)
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
		c.Break()
	}
}

// WrapHandleFunc wrap for context handler chain
func WrapHandlerFunc(h HandleFunc) HandleFunc {
	return func(c *Context) {
		h(c)
		c.Next()
	}
}
