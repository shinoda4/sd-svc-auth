package auth

import (
	"github.com/shinoda4/sd-svc-auth/internal/service/entity"
)

type Service struct {
	db    entity.UserRepository
	cache entity.CacheRepository
}

func NewAuthService(db entity.UserRepository, cache entity.CacheRepository) *Service {
	return &Service{db: db, cache: cache}
}
