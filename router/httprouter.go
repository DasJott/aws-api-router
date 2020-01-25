package router

import (
	"net/http"

	"github.com/dasjott/aws-api-router/context"
)

type (
	// HTTPHandlerFunc is the function you want to implement for each HTTP route
	HTTPHandlerFunc func(*context.HTTP)

	// HTTPRouter is the routing unit that takes your routes
	HTTPRouter struct {
		baseRouter
	}
)

// GET is convenient for Add(http.MethodGet, path, HTTPHandlerFunc)
func (r *HTTPRouter) GET(path string, handle ...HTTPHandlerFunc) {
	r.Add(http.MethodGet, path, handle)
}

// PUT is convenient for Add(http.MethodPut, path, HTTPHandlerFunc)
func (r *HTTPRouter) PUT(path string, handle ...HTTPHandlerFunc) {
	r.Add(http.MethodPut, path, handle)
}

// POST is convenient for Add(http.MethodPost, path, HTTPHandlerFunc)
func (r *HTTPRouter) POST(path string, handle ...HTTPHandlerFunc) {
	r.Add(http.MethodPost, path, handle)
}

// DELETE is convenient for Add(http.MethodDelete, path, HTTPHandlerFunc)
func (r *HTTPRouter) DELETE(path string, handle ...HTTPHandlerFunc) {
	r.Add(http.MethodDelete, path, handle)
}
