package container

import (
	v1 "scam/internal/endpoint/controller/http/api/v1"
	"scam/internal/endpoint/controller/http/api/v1/auth"
	"scam/internal/endpoint/controller/http/api/v1/user"
)

func (c *Container) getHTTPDispatcher() *v1.Dispatcher {
	if c.httpDispatcher == nil {
		c.httpDispatcher = v1.NewDispatcher(
			c.getConfig().Internal.Path,

			auth.NewController(
				c.getLogger(),
				c.getResponseBuilder(),
				c.getApplication().getAuthApplicationService(),
			),

			user.NewController(
				c.getLogger(),
				c.getResponseBuilder(),
				c.getApplication().getUserApplicationService(),
			),
		)
	}
	return c.httpDispatcher
}
