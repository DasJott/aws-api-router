package api

import "github.com/dasjott/aws-api-router/router"

type (
	// M is for creating quick json responses
	M map[string]interface{}
)

// NewREST returns a pointer to a new RESTRouter object
func NewREST() *router.RESTRouter {
	r := &router.RESTRouter{}
	return r
}

// NewHTTP returns a pointer to a new HTTPRouter object
func NewHTTP() *router.HTTPRouter {
	r := &router.HTTPRouter{}
	return r
}
