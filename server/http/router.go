package http

type Router struct {
	httpServer *HttpServer
	handles    map[string][]HandleFunc
}

func NewRouter(server *HttpServer) *Router {
	r := new(Router)
	r.httpServer = server
	r.handles = make(map[string][]HandleFunc)
	return r
}

func (r *Router) Add(method string, pattern string, handles []HandleFunc) {
	r.handles[pattern] = handles
}

func (r *Router) Get(pattern string) []HandleFunc {
	return r.handles[pattern]
}
