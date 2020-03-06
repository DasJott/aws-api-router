package router

import (
	"net/http"

	"github.com/dasjott/aws-api-router/context"
)

type (
	// RESTHandlerFunc is the function you want to implement for each REST route
	RESTHandlerFunc func(*context.REST)

	// RESTGroup groups routes on one common base path
	RESTGroup struct {
		router   *RESTRouter
		basepath string
	}

	// RESTRouter is the routing unit that takes your routes
	RESTRouter struct {
		baseRouter
		RESTGroup
	}
)

// NewRESTRouter returns a pointer to a new RESTRouter object
func NewRESTRouter() *RESTRouter {
	r := &RESTRouter{}
	r.routes = make(map[string]branch)
	r.preHandler = make(map[string]interface{})
	r.RESTGroup = RESTGroup{
		router: r,
	}
	return r
}

// Group returns a new group
func (r *RESTRouter) Group(path string) *RESTGroup {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return &RESTGroup{
		router:   r,
		basepath: path,
	}
}

// Pre defines handler functions to be executed before every
func (g *RESTGroup) Pre(handle ...RESTHandlerFunc) {
	g.router.preHandler[g.basepath] = handle
}

// Add adds a new path to the router / group
// handlerFunc can be a RESTHandlerFunc, HTTPHandlerFunc (depending on wich router you use)
// it may also be a slice of multiple of those handlers.
func (g *RESTGroup) Add(method, path string, handle ...RESTHandlerFunc) {
	if path[0] != '/' {
		path = "/" + path
	}
	path = g.basepath + path
	g.router.add(method, path, g.basepath, handle)
}

// GET creates a new GET route
func (g *RESTGroup) GET(path string, handle ...RESTHandlerFunc) {
	g.Add(http.MethodGet, path, handle...)
}

// PUT creates a new PUT route
func (g *RESTGroup) PUT(path string, handle ...RESTHandlerFunc) {
	g.Add(http.MethodPut, path, handle...)
}

// POST creates a new POST route
func (g *RESTGroup) POST(path string, handle ...RESTHandlerFunc) {
	g.Add(http.MethodPost, path, handle...)
}

// DELETE creates a new DELETE route
func (g *RESTGroup) DELETE(path string, handle ...RESTHandlerFunc) {
	g.Add(http.MethodDelete, path, handle...)
}
