package container

import (
	"scam/internal/infrastructure/repository/file"
	"scam/internal/infrastructure/repository/social"
	tokensRepo "scam/internal/infrastructure/repository/tokens"
	userRepo "scam/internal/infrastructure/repository/user"
)

func (c *Container) getRepositories() *repositories {
	if c.repositories == nil {
		c.repositories = &repositories{c: c}
	}
	return c.repositories
}

type repositories struct {
	c *Container

	user   *userRepo.Repository
	social *social.Repository
	token  *tokensRepo.Repository
	file   *file.Repository
}

func (r *repositories) getUserRepository() *userRepo.Repository {
	if r.user == nil {
		r.user = userRepo.NewRepository(r.c.getDBPool())
	}
	return r.user
}

func (r *repositories) getTokenRepository() *tokensRepo.Repository {
	if r.token == nil {
		r.token = tokensRepo.NewRepository(r.c.getDBPool())
	}
	return r.token
}

func (r *repositories) getFileRepository() *file.Repository {
	if r.file == nil {
		r.file = file.NewRepository(r.c.getDBPool())
	}
	return r.file
}

func (r *repositories) getSocialRepository() *social.Repository {
	if r.social == nil {
		r.social = social.NewRepository(r.c.getDBPool())
	}
	return r.social
}
