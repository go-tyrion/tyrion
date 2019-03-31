package http

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	httpServer *HttpServer
	Req        *http.Request
	Resp       http.ResponseWriter
	handles    []HandleFunc
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.Req = r
	ctx.Resp = w
	ctx.handles = make([]HandleFunc, 0)
	return ctx
}

func (ctx *Context) Reset(w http.ResponseWriter, r *http.Request) *Context {
	ctx.Req = r
	ctx.Resp = w
	ctx.handles = make([]HandleFunc, 0)
	return ctx
}

func (ctx *Context) Run() {
	for _, handle := range ctx.handles {
		handle(ctx)
	}
}

func (ctx *Context) String(code int, text string) {
	ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Resp.WriteHeader(code)
	ctx.Resp.Write([]byte(text))
}

func (ctx *Context) Json(code int, v interface{}) {
	var body []byte
	var err error
	if body, err = json.Marshal(v); err != nil {
		ctx.httpServer.logger.Println(err)
		return
	}

	ctx.Resp.WriteHeader(code)
	ctx.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Resp.Write(body)
}
