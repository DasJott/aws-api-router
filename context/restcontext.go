package context

import (
	"encoding/json"

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
	ctx.response = &events.APIGatewayProxyResponse{}

	return ctx
}

// GetResponse returns the inner response object
func (c *REST) GetResponse() *events.APIGatewayProxyResponse {
	return c.response
}

// String is used to just repond with a string
func (c *REST) String(code int, str string) {
	c.response.StatusCode = code
	c.response.Body = str
}

// JSON sets a struct or a
func (c *REST) JSON(code int, obj interface{}) {
	data, err := json.Marshal(obj)

	if err == nil {
		c.response.StatusCode = code
		c.response.Body = string(data)
	} else {
		c.response.StatusCode = 500
		c.response.Body = err.Error()
	}
}

// Param is a URL parameters mapping
func (c *REST) Param(key string) string {
	return c.Params[key]
}

// Query is a URL parameters mapping
func (c *REST) Query(key string) string {
	return c.Queries[key]
}
