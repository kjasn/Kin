package kin

import (
	"net/http"
)

// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

// kin.New()
func New() *Engine {
	return &Engine{router: newRouter()}
}

// engine.Run()
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// serveHttp
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// construct a context first
	c := newContext(w, req)
	// handle by router with context
	engine.router.handle(c)
}

// engine.addRouter   not public
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) { // method -- GET/POST...
	engine.router.addRoute(method, pattern, handler)
}

// engine.GET()  add GET request
func (engine *Engine) GET(path string, handler HandlerFunc) {
	engine.addRouter("GET", path, handler)
}

// engine.POST()  add POST request
func (engine *Engine) POST(path string, handler HandlerFunc) {
	engine.addRouter("POST", path, handler)
}