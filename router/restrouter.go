package router

import (
	"net/http"

	"github.com/dasjott/aws-api-router/context"
)

type (
	// RESTHandlerFunc is the function you want to implement for each REST route
	RESTHandlerFunc func(*context.REST)

	// RESTRouter is the routing unit that takes your routes
	RESTRouter struct {
		baseRouter
	}
)

// GET creates a new GET route
func (r *RESTRouter) GET(path string, handle ...RESTHandlerFunc) {
	r.Add(http.MethodGet, path, handle)
}

// PUT creates a new PUT route
func (r *RESTRouter) PUT(path string, handle ...RESTHandlerFunc) {
	r.Add(http.MethodPut, path, handle)
}

// POST creates a new POST route
func (r *RESTRouter) POST(path string, handle ...RESTHandlerFunc) {
	r.Add(http.MethodPost, path, handle)
}

// DELETE creates a new DELETE route
func (r *RESTRouter) DELETE(path string, handle ...RESTHandlerFunc) {
	r.Add(http.MethodDelete, path, handle)
}
