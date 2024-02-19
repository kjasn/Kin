package kin

import (
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup	// RouterGroup insert into Engine
	router *router
	groups []*RouterGroup
	// serve as html render
	htmlTemplates *template.Template	// store all html templates
	funcMap template.FuncMap	// render func
}




// constructor
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
		// middlewares: make([]HandlerFunc, 0),
	}
	// engine.groups = make([]*RouterGroup, 0)
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Default 	use Logger & Recovery 
func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

// Run in some port
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP, evoked while server accept request
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc

	// construct a context first
	ctx := newContext(w, req)
	// handle by router with context
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			// for _, middleware := range group.middlewares {
			// 	middlewares = append(middlewares, middleware)
			// 	middleware(ctx)
			// }
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	ctx.handlers = middlewares	// store middlewares of current req
	ctx.engine = engine
	engine.router.handle(ctx)
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
func (group *RouterGroup) Use(handlers ...HandlerFunc) {
	group.middlewares = append(group.middlewares, handlers...)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(ctx *Context) {
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		// fetch filepath (relative path)
		// ctx.JSON(
		// 	http.StatusOK,
		// 	H {"filepath" : file}, 
		// )
		// render by net/http package
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

// serve static file
// parse request addr(relativePath) to get localstorage(root)
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")

	// register urlPattern 
	group.GET(urlPattern, handler)
}


func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// load html templates to Engine globally
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}