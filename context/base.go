package context

import "github.com/aws/aws-lambda-go/events"

type baseContext struct {
	Request    *events.APIGatewayProxyRequest
	attributes map[string]interface{}
	abort      bool
}

// Set reminds a value for you
func (c *baseContext) Set(key string, val interface{}) {
	c.attributes[key] = val
}

// Get recalls a value for you
func (c *baseContext) Get(key string) interface{} {
	return c.attributes[key]
}

// Abort prohibits execution of following handler functions.
func (c *baseContext) Abort() {
	c.abort = true
}

// Aborted tells whether a handler func want to abort further execution of following handler functions.
func (c *baseContext) Aborted() bool {
	return c.abort
}
