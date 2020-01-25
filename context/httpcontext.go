package context

import (
	"github.com/aws/aws-lambda-go/events"
)

// HTTP is provided to each callback
// Do not create it manually!
type HTTP struct {
	baseContext
	Request  *events.APIGatewayProxyRequest
	Response *events.APIGatewayProxyResponse
}

// NewHTTP returns a pointer to a new HTTP context object
func NewHTTP(req *events.APIGatewayProxyRequest) *HTTP {
	ctx := &HTTP{
		baseContext: baseContext{},
		Response:    &events.APIGatewayProxyResponse{},
	}
	ctx.Request = req
	ctx.attributes = make(map[string]interface{})

	return ctx
}
