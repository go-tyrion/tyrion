package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type Context struct {
	httpServer *HttpServer
	Req        *http.Request
	Resp       http.ResponseWriter
	handles    []HandleFunc
	step       int
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
	ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Resp.WriteHeader(code)
	ctx.Resp.Write([]byte(text))
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

	ctx.Resp.WriteHeader(code)
	ctx.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Resp.Write(body)
}

func (ctx *Context) OkJSON(v interface{}) {
	ctx.JSON(http.StatusOK, v)
}

// 获取 Header
func (ctx *Context) GetHeader(key string) string {
	return ctx.Req.Header.Get(key)
}

// 设置 Header
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Resp.Header().Set(key, value)
}

func (ctx *Context) Get(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Post(key string) interface{} {
	return ctx.Req.PostForm.Get(key)
}
