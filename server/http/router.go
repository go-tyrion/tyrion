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
	HttpStatus = map[int]string{
		200: "OK",
		302: "Found",
		304: "Not modified",
		400: "Bad request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not found",
		405: "Method not allowed",
		406: "Not acceptable",
		408: "Request timeout",
		409: "Conflict",
		500: "Internal server error",
		502: "Bad gateway",
		503: "Service unavailable",
		504: "Gateway timeout",
		598: "Network read timeout error",
		599: "Network connect timeout error",
	}
)

type Router struct {
	httpServer *HttpService
	handles    map[int]map[string][]HandleFunc
}

func newRouter(server *HttpService) *Router {
	r := new(Router)
	r.httpServer = server
	r.handles = make(map[int]map[string][]HandleFunc)
	for _, method := range HttpMethods {
		r.handles[method] = make(map[string][]HandleFunc)
	}
	return r
}

func (r *Router) Register(method string, pattern string, handles []HandleFunc) {
	pattern = strings.TrimRight(strings.ToLower(pattern), "/")
	r.handles[HttpMethods[method]][pattern] = handles
}

func (r *Router) Get(method string, pattern string) []HandleFunc {
	pattern = strings.ToLower(pattern)
	if r.httpServer.opts.IgnorePathLastSlash {
		pattern = strings.TrimRight(pattern, "/")
	}

	if handles, exists := r.handles[HttpMethods[method]][pattern]; exists {
		return handles
	}

	return nil
}
