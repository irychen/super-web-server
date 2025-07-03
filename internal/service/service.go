package service

import (
	"super-web-server/internal/repo"
	"super-web-server/pkg/jwt"
	"super-web-server/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	User() UserService
}

type service struct {
	userService UserService
	logger      *logger.Logger
	redis       *redis.Client
	jwt         *jwt.JWT
}

func NewService(repo repo.Repo, logger *logger.Logger, redis *redis.Client, jwt *jwt.JWT) Service {
	logger.Info("NewService initialized successfully")
	return &service{
		userService: NewUserService(repo.User(), logger, redis, jwt),
		logger:      logger,
		redis:       redis,
		jwt:         jwt,
	}
}

func (s *service) User() UserService {
	return s.userService
}
