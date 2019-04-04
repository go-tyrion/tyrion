package http

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	HeaderApplicationJsonCharsetUTF8 = "application/json; charset=utf-8"
	HeaderTextHtmlCharsetUTF8        = "text/html; charset=utf-8"
)

type Context struct {
	httpServer *HttpServer
	req        *http.Request
	resp       http.ResponseWriter
	handles    []HandleFunc
	step       int
}

func NewContext(w http.ResponseWriter, r *http.Request, app *HttpServer) *Context {
	ctx := new(Context)
	ctx.httpServer = app
	ctx.req = r
	ctx.resp = w
	ctx.handles = make([]HandleFunc, 0)
	return ctx
}

func (ctx *Context) reset(w http.ResponseWriter, r *http.Request) *Context {
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
	ctx.resp.WriteHeader(code)
	ctx.resp.Header().Set("Content-Type", HeaderTextHtmlCharsetUTF8)
	_, err := ctx.resp.Write([]byte(text))
	ctx.Error(err)
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
	ctx.resp.Header().Set("Content-Type", HeaderApplicationJsonCharsetUTF8)
	_, err = ctx.resp.Write(body)
	ctx.Error(err)
}

func (ctx *Context) PostArray(key string) ([]string, bool) {
	req := ctx.req
	if err := req.ParseMultipartForm(ctx.httpServer.GetMaxPostMemory()); err != nil {
		if err != http.ErrNotMultipart {
			ctx.Error(err)
		}
	}
	if values := req.PostForm[key]; len(values) > 0 {
		return values, true
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return values, true
		}
	}
	return []string{}, false
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

func (ctx *Context) Post(key string) string {
	if values, exists := ctx.PostArray(key); exists {
		return values[0]
	}
	return ""
}

func (ctx *Context) Error(err error) {
	if err == nil {
		return
	}
	ctx.Log().Println("Err:", err)
}
