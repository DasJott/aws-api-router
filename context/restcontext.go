package context

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// REST is provided to each callback
// Do not create it manually!
type REST struct {
	baseContext
	Params  map[string]string
	Queries map[string]string
	Body    string
	Caller  *events.APIGatewayProxyRequestContext

	response *events.APIGatewayProxyResponse
}

// NewREST returns a pointer to a new Context object
func NewREST(req *events.APIGatewayProxyRequest) *REST {
	ctx := &REST{
		baseContext: baseContext{},
		Body:        req.Body,
		Caller:      &req.RequestContext,
	}
	ctx.Request = req
	ctx.Params = req.PathParameters
	ctx.Queries = req.QueryStringParameters
	ctx.attributes = make(map[string]interface{})
	ctx.response = &events.APIGatewayProxyResponse{
		Headers: make(map[string]string),
	}

	return ctx
}

// GetResponse returns the inner response object
func (c *REST) GetResponse() *events.APIGatewayProxyResponse {
	return c.response
}

// SetHeader sets a header field in the response
func (c *REST) SetHeader(key, val string) {
	c.response.Headers[key] = val
}

// String is used to just repond with a string
func (c *REST) String(code int, str string) {
	c.response.StatusCode = code
	c.response.Body = str

	c.response.Headers["Content-Type"] = "application/json"
	c.response.Headers["Content-Length"] = strconv.Itoa(len(str))
}

// JSON sets a struct or a
func (c *REST) JSON(code int, obj interface{}) {
	data, err := json.Marshal(obj)

	if err != nil {
		c.serverError(err)
	} else {
		c.response.StatusCode = code
		c.response.Body = string(data)
		c.response.Headers["Content-Type"] = "application/json; charset=utf-8"
		c.response.Headers["Content-Length"] = strconv.Itoa(len(c.response.Body))
	}
}

func (c *REST) serverError(err error) {
	c.response.StatusCode = 500
	c.response.Body = err.Error()
	c.response.Headers["Content-Type"] = "text/plain; charset=utf-8"
	c.response.Headers["Content-Length"] = strconv.Itoa(len(c.response.Body))
}

// Param is a URL parameters mapping
func (c *REST) Param(key string) string {
	return c.Params[key]
}

// Query is a URL parameters mapping
func (c *REST) Query(key string) string {
	return c.Queries[key]
}
