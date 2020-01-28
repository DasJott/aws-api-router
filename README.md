# AWS-API-ROUTER
A simple, small and fast router for APIs using Lambda and APIGateway from AWS<br>
Written in GO it helps to easily develop quick APIs.<br>
<br>
This library contains a REST optimized API and a plain HTTP API router.<br>
[_Me Like!_](paypal.me/dasjott)

## Usage
This library is intended to be used on a Lambda environment, but you can also use it as long as you have a APIGatewayProxyRequest object available and need back a APIGatewayProxyResponse.

Example code for an API on a Lambda:
``` go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	api "github.com/dasjott/aws-api-router"
	"github.com/dasjott/aws-api-router/context"
)

func main() {
	r := api.NewREST()

	r.GET("/myservice/{user}", func(c *context.REST) {
		c.String(200, "Hello "+c.Param("user"))
	})

	lambda.Start(r.Handle)
}
```

As the path parameter implies, for registration of a new route simply copy the path from your API Gateway on AWS.

## License
Free to use, wherever you want. Don't claim it as yours and do not sell it.<br>
If you improve or enhance it, provide it as a pull request, so everyone can enjoy.

## Donate?
If you appreciate what I do, please consider [donation via PayPal](paypal.me/dasjott).<br>
Thank you ;-)
