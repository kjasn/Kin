package kin

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

// kin.New()
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// engine.Run()
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// serveHttp
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok:= engine.router[req.Method + "-" + req.URL.Path]; ok{
		handler(w, req)
	} else {
		// not exists
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

// engine.addRouter   not public
func (e *Engine) addRouter(method string, pattern string, handler HandlerFunc) { // method -- GET/POST...
	key := method + "-" + pattern
	e.router[key] = handler
}

// engine.GET()  add GET request
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRouter("GET", path, handler)
}

// engine.POST()  add POST request
func (e *Engine) POST(path string, handler HandlerFunc) {
	e.addRouter("POST", path, handler)
}