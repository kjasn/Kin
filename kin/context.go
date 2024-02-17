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
	// kv pair storage for request
	Keys map[string]interface{}
	// engine pointer
	engine *Engine
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

func (ctx *Context) Fail(code int, err string) {
	ctx.index = len(ctx.handlers)
	ctx.JSON(
		code, 
		H{"message" : err},
	)
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

func (ctx *Context) HTML(code int, fileName string, data interface{}) {
	ctx.Writer.Header().Set("Content-Type", "text/html")
	// ctx.Status(code)
	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.Writer, fileName, data); err != nil {
		ctx.Fail(500, err.Error())
	} else {
		ctx.Status(code)
	}
	// ctx.Writer.Write([]byte(html))
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.Writer.Write(data)
}

func (ctx *Context) Next() {
	ctx.index ++	// current middleware index
	n := len(ctx.handlers)

	for ; ctx.index < n; ctx.index ++ {	// exec in order
		ctx.handlers[ctx.index](ctx)
	}
}

func (ctx *Context) Set(key string, value interface{}) {
	if ctx.Keys == nil {
		ctx.Keys = make(map[string]interface{})
	}

	ctx.Keys[key] = value
}


func (ctx *Context) Get(key string) (value interface{}, ok bool) {
	value, ok = ctx.Keys[key]
	return
}