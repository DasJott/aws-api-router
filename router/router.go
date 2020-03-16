package router

import (
	"fmt"
	"reflect"
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
		Func  interface{}
		Group string
	}

	// baseRouter is the base for every router
	baseRouter struct {
		routes     map[string]branch
		preHandler map[string]interface{}
	}
	// HTTPErrorHandler handles occuring rest errors
	HTTPErrorHandler func(error, *context.HTTP)
	// RESTErrorHandler handles occuring rest errors
	RESTErrorHandler func(error, *context.REST)
)

// ErrorHandler must be one of HTTPErrorHandler or RESTErrorHandler and can be set to handle emerging errors your way
var ErrorHandler interface{}

// add adds a new path to the router
// handlerFunc can be a RESTHandlerFunc, HTTPHandlerFunc (depending on wich router you use)
// it may also be a slice of multiple of those handlers.
func (r *baseRouter) add(method, path, group string, handlerFunc interface{}) {
	if handlerFunc == nil {
		return
	}

	path = strings.Trim(path, "/")

	h := &handler{
		Func:  handlerFunc,
		Group: group,
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
			if m[part].Branch == nil {
				m[part].Branch = make(branch)
			}
			m = m[part].Branch
		}
	}
}

// Find finds the handler to the given path
func (r *baseRouter) find(req *events.APIGatewayProxyRequest) *handler {
	parts := strings.Split(req.Resource[1:], "/")

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
		list := r.getHandlerFunctionList(h)

		switch funcs := list.(type) {
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

	err := fmt.Errorf("route not found:" + req.Resource)

	if ErrorHandler != nil {
		switch handler := ErrorHandler.(type) {
		case RESTErrorHandler:
			c := context.NewREST(req)
			handler(err, c)
			return c.GetResponse(), err
		case HTTPErrorHandler:
			c := context.NewHTTP(req)
			handler(err, c)
			return c.Response, err
		}
	}

	return &events.APIGatewayProxyResponse{
		Body:       err.Error(),
		StatusCode: 404,
	}, err
}

func (r *baseRouter) getHandlerFunctionList(h *handler) interface{} {
	arr := reflect.MakeSlice(reflect.TypeOf(h.Func), 0, 10)
	if g := r.preHandler[h.Group]; g != nil {
		arr = reflect.AppendSlice(arr, reflect.ValueOf(g))
	}
	arr = reflect.AppendSlice(arr, reflect.ValueOf(h.Func))
	return arr.Interface()
}
