package api

import (
	"github.com/dasjott/aws-api-router/router"
)

type (
	// M is for creating quick json responses
	M map[string]interface{}
)

// NewREST returns a pointer to a new RESTRouter object
func NewREST() *router.RESTRouter {
	return router.NewREST()
}

// NewHTTP returns a pointer to a new HTTPRouter object
func NewHTTP() *router.HTTPRouter {
	return router.NewHTTPRouter()
}
