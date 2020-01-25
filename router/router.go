package router

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dasjott/aws-api-router/context"
)

type (
	node struct {
		Branch  branch
		Handler *handler
	}
	branch map[string]*node

	handler struct {
		Func interface{}
	}

	// baseRouter is the base for every router
	baseRouter struct {
		// http.method -> map[string]*Route or anther map[string]interface{}
		routes map[string]branch
	}
)

// NewREST returns a pointer to a new RESTRouter object
func NewREST() *RESTRouter {
	r := &RESTRouter{}
	r.routes = make(map[string]branch)
	return r
}

// NewHTTP returns a pointer to a new HTTPRouter object
func NewHTTP() *HTTPRouter {
	r := &HTTPRouter{}
	r.routes = make(map[string]branch)
	return r
}

// Add adds a new path to the router
func (r *baseRouter) Add(method, path string, handlerFunc interface{}) {
	path = strings.Trim(path, "/")

	h := &handler{
		Func: handlerFunc,
	}

	pathParts := strings.Split(path, "/")
	last := len(pathParts) - 1
	if last > 255 {
		panic("path too long: " + path)
	}

	if r.routes[method] == nil {
		r.routes[method] = make(branch)
	}
	m := r.routes[method]

	for i, part := range pathParts {
		if partlen := len(part) - 1; partlen > 1 && part[0] == '{' && part[partlen] == '}' {
			part = "*"
		}

		if m[part] == nil {
			m[part] = &node{}
		}
		if i == last {
			m[part].Handler = h
		} else {
			m[part].Branch = make(branch)
			m = m[part].Branch
		}
	}
}

// Find finds the handler to the given path
func (r *baseRouter) find(req *events.APIGatewayProxyRequest) *handler {
	parts := strings.Split(req.Path[1:], "/")

	if m, ok := r.routes[req.HTTPMethod]; ok {
		last := len(parts) - 1
		for i, part := range parts {
			if m[part] == nil {
				part = "*"
			}
			if m[part] == nil {
				break
			}
			if i == last {
				return m[part].Handler
			}
			m = m[part].Branch
		}
	}

	return nil
}

// Handle is the function to handle API Gateway requests.
// Put this function into the lambda.Start function.
func (r *baseRouter) Handle(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if h := r.find(req); h != nil {
		switch funcs := h.Func.(type) {
		// Single handler
		case func(*context.REST):
			c := context.NewREST(req)
			funcs(c)
			return c.GetResponse(), nil
		case func(*context.HTTP):
			c := context.NewHTTP(req)
			funcs(c)
			return c.Response, nil

		// Multihandler
		case []RESTHandlerFunc:
			c := context.NewREST(req)
			for _, f := range funcs {
				f(c)
			}
			return c.GetResponse(), nil
		case []HTTPHandlerFunc:
			c := context.NewHTTP(req)
			for _, f := range funcs {
				f(c)
			}
			return c.Response, nil
		}
	}

	return nil, fmt.Errorf("route not found")
}
