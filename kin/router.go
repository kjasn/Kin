package kin

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

// define router group
type RouterGroup struct {
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup	
	engine *Engine
}

// not public
func newRouter() *router {
	return &router{
		roots: make(map[string]*node),	// kv: request method -- router path  e.g. roots["GET"]=/
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {	// root path not exist
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}


func parsePattern(path string) []string {
	ps := strings.Split(path, "/")

	parts := make([]string, 0)
	for _, itr := range ps {
		if itr != "" {
			parts = append(parts, itr)
			if itr[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) 
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0, false)	// parsed path
	if node != nil {
		parts := parsePattern(node.pattern)

		for idx, part := range parts {
			if part[0] == ':' { // ':' parameter, /:lang --- /go => lang = go
				params[part[1 : ]] = searchParts[idx]
			}

			if part[0] == '*' && len(parts) > 1 {
				// '*' -- wildcard, join all remain parts as complete path
				params[part[1 : ]] = strings.Join(searchParts[idx : ], "/")	
				break
			}
		}
		return node, params
	}

	return nil, nil
}

func (r *router) handle(ctx *Context) {
	n, params := r.getRoute(ctx.Method, ctx.Path)
	
	if n != nil {
		ctx.Params = params
		key := ctx.Method + "-" + n.pattern
		ctx.handlers = append(ctx.handlers, r.handlers[key])
		// r.handlers[key](ctx)
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		})
	}

	ctx.Next()
}