package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type Context struct {
	httpServer *HttpServer
	req        *http.Request
	resp       http.ResponseWriter
	handles    []HandleFunc
	step       int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.req = r
	ctx.resp = w
	ctx.handles = make([]HandleFunc, 0)
	return ctx
}

func (ctx *Context) Reset(w http.ResponseWriter, r *http.Request) *Context {
	ctx.req = r
	ctx.resp = w
	ctx.handles = make([]HandleFunc, 0)
	ctx.step = 0
	return ctx
}

func (ctx *Context) Next() {
	if ctx.step >= len(ctx.handles) {
		return
	}

	i := ctx.step
	ctx.step++

	ctx.handles[i](ctx)
}

func (ctx *Context) Log() *log.Logger {
	return ctx.httpServer.Log()
}

func (ctx *Context) Break() {
	ctx.step = len(ctx.handles)
}

func (ctx *Context) String(code int, text string) {
	ctx.resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.resp.WriteHeader(code)
	ctx.resp.Write([]byte(text))
}

func (ctx *Context) OkString(text string) {
	ctx.String(http.StatusOK, text)
}

func (ctx *Context) JSON(code int, v interface{}) {
	var body []byte
	var err error
	if body, err = json.Marshal(v); err != nil {
		ctx.httpServer.logger.Println(err)
		return
	}

	ctx.resp.WriteHeader(code)
	ctx.resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.resp.Write(body)
}

func (ctx *Context) OkJSON(v interface{}) {
	ctx.JSON(http.StatusOK, v)
}

// 获取 Header
func (ctx *Context) GetHeader(key string) string {
	return ctx.req.Header.Get(key)
}

// 设置 Header
func (ctx *Context) SetHeader(key string, value string) {
	ctx.resp.Header().Set(key, value)
}

func (ctx *Context) Get(key string) string {
	return ctx.req.URL.Query().Get(key)
}

func (ctx *Context) Post(key string) interface{} {
	return ctx.req.PostForm.Get(key)
}
