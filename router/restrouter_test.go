package router_test

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dasjott/aws-api-router/context"
	"github.com/dasjott/aws-api-router/router"
	"github.com/stretchr/testify/assert"
)

func TestMultihandler(t *testing.T) {
	test := assert.New(t)
	r := router.NewREST()

	r.GET("/foo/bar/baz/moinsen",
		func(c *context.REST) { c.Set("memory", "drink?") },
		func(c *context.REST) { c.String(200, c.Get("memory").(string)+" coffee!") },
	)

	{ // find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz/moinsen",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		ok := test.Nil(err)
		ok = ok && test.NotNil(resp)

		ok = ok && test.EqualValues("drink? coffee!", resp.Body)
		ok = ok && test.EqualValues(200, resp.StatusCode)
	}
}

func TestFind(t *testing.T) {
	test := assert.New(t)
	r := router.NewREST()

	r.GET("/foo/bar/baz/moin", func(c *context.REST) { c.String(418, "I'm a teapot") })

	{ // find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz/moin",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.Nil(err)
		test.NotNil(resp)

		test.EqualValues("I'm a teapot", resp.Body)
		test.EqualValues(418, resp.StatusCode)
	}

	{ // don't find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.NotNil(err)
		test.Nil(resp)
	}

	{ // don't find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz/bummer",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.NotNil(err)
		test.Nil(resp)
	}

	{ // don't find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz/moin/bummer",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.NotNil(err)
		test.Nil(resp)
	}

	r.GET("/foo/bar/baz/moinsen", func(c *context.REST) { c.String(419, "coffee?") })

	{ // don't find the other
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz/moinsen",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		ok := test.Nil(err)
		ok = ok && test.NotNil(resp)

		ok = ok && test.EqualValues("coffee?", resp.Body)
		ok = ok && test.EqualValues(419, resp.StatusCode)
	}
}

func TestFindWithParams(t *testing.T) {
	test := assert.New(t)
	r := router.NewREST()

	r.GET("/foo/{b}/moin/{name}", func(c *context.REST) {
		c.String(200, c.Param("b")+" "+c.Param("name"))
	})

	{ // find it
		req := &events.APIGatewayProxyRequest{
			Path:           "/foo/bar/moin/jott",
			HTTPMethod:     http.MethodGet,
			PathParameters: map[string]string{"b": "bar", "name": "jott"},
		}
		resp, err := r.Handle(req)
		test.Nil(err)
		test.NotNil(resp)

		test.EqualValues("bar jott", resp.Body)
		test.EqualValues(200, resp.StatusCode)
	}

	{ // don't find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.Nil(resp)
		test.NotNil(err)
	}

	{ // don't find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz/bummer",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.Nil(resp)
		test.NotNil(err)
	}

	{ // don't find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/moin/bummer/yeah",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.Nil(resp)
		test.NotNil(err)
	}

	r.GET("/foo/bar/baz/moinsen", func(c *context.REST) {
		if len(c.Params) > 0 {
			c.String(400, "wrong")
		} else {
			c.String(202, "right")
		}
	})

	{ // don't find the other
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/baz/moinsen",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.Nil(err)
		test.NotNil(resp)

		test.EqualValues("right", resp.Body)
		test.EqualValues(202, resp.StatusCode)
	}

	// similar but no param
	r.GET("/foo/bar/moin/moinsen", func(c *context.REST) {
		if len(c.Params) > 0 {
			c.String(400, "wrong")
		} else {
			c.String(202, "right")
		}
	})

	{ // don't find the other
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar/moin/moinsen",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		test.Nil(err)
		test.NotNil(resp)

		test.EqualValues("right", resp.Body)
		test.EqualValues(202, resp.StatusCode)
	}
}

func TestGroups(t *testing.T) {
	test := assert.New(t)
	r := router.NewREST()
	r.Add(http.MethodGet, "/find/me", func(c *context.REST) { c.String(200, "direct") })

	g := r.Group("/basic")
	g.GET("/find/me", func(c *context.REST) { c.String(200, "group") })

	{ // find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/find/me",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		ok := test.Nil(err)
		ok = ok && test.NotNil(resp)

		ok = ok && test.EqualValues("direct", resp.Body)
		ok = ok && test.EqualValues(200, resp.StatusCode)
	}
	{ // find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/basic/find/me",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		ok := test.Nil(err)
		ok = ok && test.NotNil(resp)

		ok = ok && test.EqualValues("group", resp.Body)
		ok = ok && test.EqualValues(200, resp.StatusCode)
	}

}

func TestPreHandler(t *testing.T) {
	test := assert.New(t)
	r := router.NewREST()

	r.Pre(func(c *context.REST) { c.Set("value", "moin") })
	r.GET("/foo/bar", func(c *context.REST) {
		c.String(200, c.Get("value").(string)+" people")
	})

	{ // find it
		req := &events.APIGatewayProxyRequest{
			Path:       "/foo/bar",
			HTTPMethod: http.MethodGet,
		}
		resp, err := r.Handle(req)
		ok := test.Nil(err)
		ok = ok && test.NotNil(resp)

		ok = ok && test.EqualValues("moin people", resp.Body)
		ok = ok && test.EqualValues(200, resp.StatusCode)
	}

}
