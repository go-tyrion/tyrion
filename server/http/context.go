package http

import (
	"encoding/json"
	"lib/log"
	"net"
	"net/http"
	"strings"
)

const (
	HeaderApplicationJsonCharsetUTF8 = "application/json; charset=utf-8"
	HeaderTextHtmlCharsetUTF8        = "text/html; charset=utf-8"
)

type Context struct {
	httpServer *HttpService
	req        *http.Request
	resp       http.ResponseWriter
	handles    []HandleFunc
	step       int
}

func newContext(server *HttpService) *Context {
	return &Context{
		httpServer: server,
		handles:    make([]HandleFunc, 0),
		step:       0,
	}
}

func (c *Context) reset(w http.ResponseWriter, r *http.Request) *Context {
	c.req = r
	c.resp = w
	c.handles = make([]HandleFunc, 0)
	c.step = 0
	return c
}

func (c *Context) handleHTTPRequest() {
	var status int
	defer func() {
		if c.httpServer.opts.AccessLog {
			c.httpServer.accessLogger.Print(
				c.req.Method,
				status,
				c.req.URL.RequestURI(),
				c.IP())
		}
	}()

	if _, ok := HttpMethods[c.req.Method]; !ok {
		status = 405
		c.handles = append(c.handles, catchHandles(status))
		return
	}

	handles := c.httpServer.router.Get(c.req.Method, c.req.URL.Path)
	if handles == nil {
		status = 404
		c.handles = append(c.handles, catchHandles(status))
	} else {
		status = 200
		c.handles = handles
	}

	c.Next()
}

func (c *Context) Next() {
	if c.step >= len(c.handles) {
		return
	}

	i := c.step
	c.step++

	c.handles[i](c)
}

func (c *Context) Log() *log.Logger {
	return c.httpServer.Log()
}

func (c *Context) Break() {
	c.step = len(c.handles)
}

func (c *Context) IP() string {
	if true {
		clientIP := c.GetHeader("X-Forwarded-For")
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
		if clientIP == "" {
			clientIP = strings.TrimSpace(c.GetHeader("X-Real-Ip"))
		}
		if clientIP != "" {
			return clientIP
		}
	}

	if true {
		if addr := c.GetHeader("X-Appengine-Remote-Addr"); addr != "" {
			return addr
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.req.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func (c *Context) String(code int, text string) {
	c.resp.WriteHeader(code)
	c.resp.Header().Set("Content-Type", HeaderTextHtmlCharsetUTF8)
	_, err := c.resp.Write([]byte(text))
	c.Error(err)
}

func (c *Context) OkString(text string) {
	c.String(http.StatusOK, text)
}

func (c *Context) JSON(code int, v interface{}) {
	var body []byte
	var err error
	if body, err = json.Marshal(v); err != nil {
		c.httpServer.logger.Error(err)
		return
	}

	c.resp.WriteHeader(code)
	c.resp.Header().Set("Content-Type", HeaderApplicationJsonCharsetUTF8)
	_, err = c.resp.Write(body)
	c.Error(err)
}

func (c *Context) PostArray(key string) ([]string, bool) {
	req := c.req
	if err := req.ParseMultipartForm(c.httpServer.opts.GetMaxPostMemory()); err != nil {
		if err != http.ErrNotMultipart {
			c.Error(err)
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

func (c *Context) OkJSON(v interface{}) {
	c.JSON(http.StatusOK, v)
}

// 获取 Header
func (c *Context) GetHeader(key string) string {
	return c.req.Header.Get(key)
}

// 设置 Header
func (c *Context) SetHeader(key string, value string) {
	c.resp.Header().Set(key, value)
}

func (c *Context) Get(key string) string {
	return c.req.URL.Query().Get(key)
}

func (c *Context) Post(key string) string {
	if values, exists := c.PostArray(key); exists {
		return values[0]
	}
	return ""
}

func (c *Context) Error(err error) {
	if err == nil {
		return
	}
	c.Log().Error("Err:", err)
}
