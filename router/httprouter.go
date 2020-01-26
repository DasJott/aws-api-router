package router

import (
	"net/http"

	"github.com/dasjott/aws-api-router/context"
)

type (
	// HTTPHandlerFunc is the function you want to implement for each HTTP route
	HTTPHandlerFunc func(*context.HTTP)

	// HTTPGroup groups routes on one common base path
	HTTPGroup struct {
		router   *HTTPRouter
		basepath string
	}

	// HTTPRouter is the routing unit that takes your routes
	HTTPRouter struct {
		baseRouter
		HTTPGroup
	}
)

// NewHTTPRouter returns a pointer to a new HTTPRouter object
func NewHTTPRouter() *HTTPRouter {
	r := &HTTPRouter{}
	r.HTTPGroup = HTTPGroup{
		router: r,
	}
	return r
}

// Group returns a new group
func (r *HTTPRouter) Group(path string) *HTTPGroup {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return &HTTPGroup{
		router:   r,
		basepath: path,
	}
}

// Add adds a new path to the router / group
// handlerFunc can be a RESTHandlerFunc, HTTPHandlerFunc (depending on wich router you use)
// it may also be a slice of multiple of those handlers.
func (g *HTTPGroup) Add(method, path string, handle ...HTTPHandlerFunc) {
	if path[0] != '/' {
		path = "/" + path
	}
	path = g.basepath + path
	g.router.add(method, path, handle)
}

// GET is convenient for Add(http.MethodGet, path, HTTPHandlerFunc)
func (g *HTTPGroup) GET(path string, handle ...HTTPHandlerFunc) {
	g.Add(http.MethodGet, path, handle...)
}

// PUT is convenient for Add(http.MethodPut, path, HTTPHandlerFunc)
func (g *HTTPGroup) PUT(path string, handle ...HTTPHandlerFunc) {
	g.Add(http.MethodPut, path, handle...)
}

// POST is convenient for Add(http.MethodPost, path, HTTPHandlerFunc)
func (g *HTTPGroup) POST(path string, handle ...HTTPHandlerFunc) {
	g.Add(http.MethodPost, path, handle...)
}

// DELETE is convenient for Add(http.MethodDelete, path, HTTPHandlerFunc)
func (g *HTTPGroup) DELETE(path string, handle ...HTTPHandlerFunc) {
	g.Add(http.MethodDelete, path, handle...)
}
