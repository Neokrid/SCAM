package container

import (
	"scam/internal/domain/services/file"
	smtpSrv "scam/internal/domain/services/smtp"
	"scam/internal/domain/services/social"
	tokenSrv "scam/internal/domain/services/token"
	userSrv "scam/internal/domain/services/user"
)

func (c *Container) getServices() *services {
	if c.services == nil {
		c.services = &services{c: c}
	}
	return c.services
}

type services struct {
	c *Container

	user   *userSrv.Service
	social *social.Service
	smtp   *smtpSrv.Service
	token  *tokenSrv.Service
	file   *file.Service
}

func (s *services) getUserService() *userSrv.Service {
	if s.user == nil {
		s.user = userSrv.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),
			s.c.getRepositories().getUserRepository(),
		)
	}
	return s.user
}

func (s *services) getSMTPService() *smtpSrv.Service {
	if s.smtp == nil {
		s.smtp = smtpSrv.NewService(
			s.c.getLogger(),
			smtpSrv.NewConfig(
				s.c.getConfig().Email.OwnerEmail,
				s.c.getConfig().Email.OwnerPassword,
				s.c.getConfig().Email.Address,
				s.c.getConfig().Email.CodeLenght,
				s.c.getConfig().Email.CodeExp,
				s.c.getConfig().Email.MinTTL,
			),
			s.c.getCaches().getSmtpCache(),
		)
	}
	return s.smtp
}

func (s *services) getTokenService() *tokenSrv.Service {
	if s.token == nil {
		s.token = tokenSrv.NewService(
			s.c.getConfig().Jwt.AccessTTL,
			s.c.getConfig().Jwt.RefreshTTL,
			s.c.getConfig().Jwt.JwtSecret,
			s.c.getRepositories().getTokenRepository(),
		)

	}
	return s.token
}

func (s *services) getFileService() *file.Service {
	if s.file == nil {
		s.file = file.NewService(
			s.c.getLogger(),
			s.c.getRepositories().getFileRepository(),
		)

	}
	return s.file
}

func (s *services) getSocialService() *social.Service {
	if s.social == nil {
		s.social = social.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),
			s.c.getRepositories().getSocialRepository(),
		)
	}
	return s.social
}
