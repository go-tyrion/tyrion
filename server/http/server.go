package http

import (
	"fmt"
	"lib/config"
	"lib/core"
	"lib/log"
	"net/http"
	"reflect"
	"sync"
)

type HandleFunc func(c *Context)

type HttpService struct {
	core.App

	opts   *Options
	router *Router
	server *http.Server
	logger *log.Logger
	pool   sync.Pool
}

func NewHttpService() *HttpService {
	service := &HttpService{
		logger: log.NewLogger(),
		server: new(http.Server),
		opts:   new(Options),
	}
	service.router = newRouter(service)
	service.pool.New = func() interface{} {
		return newContext(service)
	}

	service.App.Init()

	err := config.Resolve("http", &service.opts.HttpConfig)
	if err != nil {
		fmt.Println("err:", err)
	}

	fmt.Println("conf", service.opts.HttpConfig)

	return service
}

// 通过 Init 方法初始化
func (s *HttpService) Init(opts *Options) {
	s.opts = s.opts.ResetOpts(opts)
	s.server.Addr = opts.Addr
	s.server.ReadTimeout = s.opts.ReadTimeout
	s.server.WriteTimeout = s.opts.WriteTimeout
}

// 通过配置文件初始化
func (s *HttpService) InitByConfig(confFile string) {
	s.Init(s.opts.ResolveOptsByConfigFile(confFile))
}

func (s *HttpService) Log() *log.Logger {
	return s.logger
}

// Run http server
func (service *HttpService) Run() error {
	service.setServerOpts()
	return service.server.ListenAndServe()
}

// Run https server
func (service *HttpService) RunTLS() error {
	if service.opts.HttpsCertFile == "" || service.opts.HttpsKeyFile == "" {
		panic("invalid tls config")
	}

	service.setServerOpts()
	return service.server.ListenAndServeTLS(service.opts.HttpsCertFile, service.opts.HttpsKeyFile)
}

func (service *HttpService) setServerOpts() {
	service.server.WriteTimeout = service.opts.WriteTimeout
	service.server.ReadTimeout = service.opts.ReadTimeout
	service.server.Handler = service
	// server.server.TLSConfig = tls.NewConfig("", "")
}

// ------------
// 类型 bisinessServer模式
func (s *HttpService) AddLogic(prefix string, logic Logic) {
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

func (s *HttpService) wrapLogic(v reflect.Value) HandleFunc {
	return func(c *Context) {
		v.Call([]reflect.Value{reflect.ValueOf(c)})
		c.Next()
	}
}

// 传统方式
func (s *HttpService) Any(pattern string, h ...HandleFunc) {
	s.add(http.MethodGet, pattern, h)
	s.add(http.MethodPost, pattern, h)
	s.add(http.MethodPut, pattern, h)
	s.add(http.MethodDelete, pattern, h)
}

func (s *HttpService) Get(pattern string, h ...HandleFunc) {
	s.add(http.MethodGet, pattern, h)
}

func (s *HttpService) Post(pattern string, h ...HandleFunc) {
	s.add(http.MethodPost, pattern, h)
}

func (s *HttpService) Put(pattern string, h ...HandleFunc) {
	s.add(http.MethodPut, pattern, h)
}

func (s *HttpService) Delete(pattern string, h ...HandleFunc) {
	s.add(http.MethodDelete, pattern, h)
}

func (s *HttpService) add(method string, pattern string, handles []HandleFunc) {
	wrapHandles := make([]HandleFunc, 0, len(handles))
	for _, h := range handles {
		wrapHandles = append(wrapHandles, WrapHandlerFunc(h))
	}
	s.router.Register(method, pattern, wrapHandles)
}

func (s *HttpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (s *HttpService) Stop() error {
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
