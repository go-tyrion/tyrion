package http

import (
	"strings"
)

const (
	GET = iota
	HEAD
	POST
	PUT
	PATCH
	DELETE
	CONNECT
	OPTIONS
	TRACE
)

var (
	HttpMethods = map[string]int{
		"GET":     GET,
		"HEAD":    HEAD,
		"POST":    POST,
		"PUT":     PUT,
		"PATCH":   PATCH,
		"DELETE":  DELETE,
		"CONNECT": CONNECT,
		"OPTIONS": OPTIONS,
		"TRACE":   TRACE,
	}
)

type Router struct {
	httpServer *HttpServer
	handles    map[int]map[string][]HandleFunc
}

func NewRouter(server *HttpServer) *Router {
	r := new(Router)
	r.httpServer = server
	r.handles = make(map[int]map[string][]HandleFunc)
	for _, method := range HttpMethods {
		r.handles[method] = make(map[string][]HandleFunc)
	}
	return r
}

func (r *Router) Register(method string, pattern string, handles []HandleFunc) {
	pattern = strings.TrimRight(pattern, "/")
	r.handles[HttpMethods[method]][pattern] = handles
}

func (r *Router) Get(method string, pattern string) []HandleFunc {
	if r.httpServer.opts.IgnorePathLastSlash {
		pattern = strings.TrimRight(pattern, "/")
	}

	if handles, exists := r.handles[HttpMethods[method]][pattern]; exists {
		return handles
	}

	return nil
}
