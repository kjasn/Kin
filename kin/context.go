package kin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index int
}

// construction
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index: -1,
	}
}

// Request
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Param(key string) string {
	// val, _ := ctx.Params[key]
	// return val
	return ctx.Params[key]
}

// Writer   Response
func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) SetHeader(key string, val string) {
	ctx.Writer.Header().Set(key, val)
}

func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.Writer.Header().Set("Content-Type", "text/plain")
	ctx.Status(code)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

func (ctx *Context) HTML(code int, html string) {
	ctx.Writer.Header().Set("Content-Type", "application/html")
	ctx.Status(code)
	ctx.Writer.Write([]byte(html))
}


func (ctx *Context) Next() {
	ctx.index ++	// current middleware index
	n := len(ctx.handlers)

	for ; ctx.index < n; ctx.index ++ {	// exec in order
		ctx.handlers[ctx.index](ctx)
	}
}

// #TODO
func (ctx *Context) Set(key string, data interface{}) {
	
}


func (ctx *Context) Get(key string) {

}