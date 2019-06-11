package http

import (
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

	opts         *Options
	router       *Router
	server       *http.Server
	logger       *log.Logger
	accessLogger *log.Logger
	pool         sync.Pool
}

func NewHttpService() *HttpService {
	service := &HttpService{
		logger:       log.NewLogger(),
		accessLogger: log.NewLogger(),
		server:       new(http.Server),
		opts:         newOptions(config.DefaultHttpConfigFile),
	}
	service.router = newRouter(service)
	service.pool.New = func() interface{} {
		return newContext(service)
	}

	service.init()

	return service
}

// 通过 Init 方法初始化
func (service *HttpService) init() {
	service.App.Init()
	service.initLog()
}

func (service *HttpService) initLog() {
	if service.opts.AccessLog {
		logger := log.NewLogger()

		if service.opts.AccessLogDir != "" {
			logger.SetOutputDir(service.opts.AccessLogDir)
			logger.SetOutputByName("access.log")

			switch service.opts.AccessLogRotate {
			case "D", "d", "day", "daily":
				logger.SetRotateDaily()
			case "H", "h", "hour", "hourly":
				logger.SetRotateHourly()
			default:
				logger.SetRotateHourly()
			}
		}

		service.accessLogger = logger
	}
}

// 通过配置文件初始化
func (s *HttpService) InitByConfig(confFile string) {

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
	service.server.Addr = service.opts.Addr
	service.server.WriteTimeout = service.opts.WriteTimeout
	service.server.ReadTimeout = service.opts.ReadTimeout
	service.server.Handler = service
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

	c.handleHTTPRequest()
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
