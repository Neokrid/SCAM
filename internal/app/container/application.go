package container

import (
	"scam/internal/application/auth"
	"scam/internal/application/social"
	userApp "scam/internal/application/user"
)

func (c *Container) getApplication() *applications {
	if c.applications == nil {
		c.applications = &applications{c: c}
	}
	return c.applications
}

type applications struct {
	c *Container

	user   *userApp.Service
	auth   *auth.Service
	social *social.Service
}

func (s *applications) getUserApplicationService() *userApp.Service {
	if s.user == nil {
		s.user = userApp.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),
			s.c.getServices().getUserService(),
			s.c.getServices().getFileService(),
			5, //добавить в env
		)
	}
	return s.user
}

func (s *applications) getAuthApplicationService() *auth.Service {
	if s.auth == nil {
		s.auth = auth.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),
			s.c.getServices().getUserService(),
			s.c.getServices().getSMTPService(),
			s.c.getServices().getTokenService(),
		)
	}
	return s.auth
}

func (s *applications) getSocialApplicationService() *social.Service {
	if s.social == nil {
		s.social = social.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),
			s.c.getServices().getSocialService(),
		)
	}
	return s.social
}
