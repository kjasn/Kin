package kin

import (
	"log"
	"net/http"
)

// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup	// RouterGroup insert into Engine
	router *router
	groups []*RouterGroup
}


// constructor
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
		middlewares: make([]HandlerFunc, 0),
	}
	engine.groups = make([]*RouterGroup, 0)
	return engine
}

// Run in some port
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP, evoked while server accept request
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// construct a context first
	c := newContext(w, req)
	// handle by router with context
	for _, group := range engine.groups {
		for _, middleware := range group.middlewares {
			middleware(c)
		}
	}

	engine.router.handle(c)
}

// not public
func (group *RouterGroup) addRouter(method string, segment string, handler HandlerFunc) { // method -- GET/POST...
	pattern := group.prefix + segment
	log.Printf("Route %s - %s\n", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// add GET request
func (group *RouterGroup) GET(path string, handler HandlerFunc) {
	group.addRouter("GET", path, handler)
}

// add POST request
func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.addRouter("POST", path, handler)
}

// create RouterGroup
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Register middleware for group 
func (group *RouterGroup) Use(handler HandlerFunc) {
	group.middlewares = append(group.middlewares, handler)
}
