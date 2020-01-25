package context

import "github.com/aws/aws-lambda-go/events"

type baseContext struct {
	Request    *events.APIGatewayProxyRequest
	attributes map[string]interface{}
}

// Set reminds a value for you
func (c *baseContext) Set(key string, val interface{}) {
	c.attributes[key] = val
}

// Get recalls a value for you
func (c *baseContext) Get(key string) interface{} {
	return c.attributes[key]
}
