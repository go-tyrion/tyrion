package http

import "net/http"

const (
	GET  = "GET"
	POST = "POST"
)

type Router struct{}

type HandFunc func(w http.ResponseWriter, r *http.Request)

func (r *Router) Add(method string, pattern string, headers ...HandFunc) {

}
