package http

type Router struct {
	httpServer *HttpServer
	handles    map[string][]HandleFunc
}

func NewRouter(app *HttpServer) *Router {
	router := new(Router)
	router.httpServer = app
	router.handles = make(map[string][]HandleFunc)
	return router
}

func (r *Router) Add(method string, pattern string, handles []HandleFunc) {
	r.handles[pattern] = handles
}

func (r *Router) Get(pattern string) []HandleFunc {
	return r.handles[pattern]
}
