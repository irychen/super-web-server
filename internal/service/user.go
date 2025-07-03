package service

import (
	"context"
	"encoding/json"
	"fmt"
	"super-web-server/internal/dto"
	"super-web-server/internal/exception"
	"super-web-server/internal/model"
	"super-web-server/internal/repo"
	"super-web-server/pkg/jwt"
	"super-web-server/pkg/logger"
	"super-web-server/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint64) (*model.User, *exception.Exception)
	GetUserByUniqueID(ctx context.Context, uniqueID int64) (*model.User, *exception.Exception)
	GetUserCachedRolesByUniqueID(ctx context.Context, uniqueID int64) ([]*model.UserRole, *exception.Exception)
	LoginByEmail(ctx context.Context, data dto.UserLoginByEmailReqDTO) (*dto.UserLoginByEmailResDTO, *exception.Exception)
}

type userService struct {
	userRepo repo.UserRepo
	logger   *logger.Logger
	redis    *redis.Client
	jwt      *jwt.JWT
}

func NewUserService(userRepo repo.UserRepo, logger *logger.Logger, redis *redis.Client, jwt *jwt.JWT) UserService {
	logger.Info("NewUserService initialized successfully")
	return &userService{
		userRepo: userRepo,
		logger:   logger,
		redis:    redis,
		jwt:      jwt,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id uint64) (*model.User, *exception.Exception) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, exception.ExceptionUserNotFound.AppendDetails(err.Error())
	}
	return user, nil
}

func (s *userService) GetUserByUniqueID(ctx context.Context, uniqueID int64) (*model.User, *exception.Exception) {
	user, err := s.userRepo.FindByUniqueID(ctx, uniqueID)
	if err != nil {
		return nil, exception.ExceptionUserNotFound.AppendDetails(err.Error())
	}
	return user, nil
}

func (s *userService) GetUserCachedRolesByUniqueID(ctx context.Context, uniqueID int64) ([]*model.UserRole, *exception.Exception) {
	cacheKey := fmt.Sprintf("user:roles:%d", uniqueID)

	cache, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil && cache != "" {
		var roles []*model.UserRole
		if err := json.Unmarshal([]byte(cache), &roles); err != nil {
			s.logger.Warn("Failed to unmarshal cached roles", zap.Int64("uniqueID", uniqueID), zap.Error(err))
		} else {
			return roles, nil
		}
	} else if err != nil && err != redis.Nil {
		s.logger.Warn("Failed to get cached roles", zap.Int64("uniqueID", uniqueID), zap.Error(err))
	}

	user, err := s.userRepo.FindByUniqueID(ctx, uniqueID)
	if err != nil {
		return nil, exception.ExceptionUserNotFound.AppendDetails(err.Error())
	}

	roles := user.Roles

	if rolesBytes, err := json.Marshal(roles); err != nil {
		s.logger.Warn("Failed to marshal roles for caching", zap.Error(err))
	} else {
		if err := s.redis.Set(ctx, cacheKey, rolesBytes, 5*time.Minute).Err(); err != nil {
			s.logger.Warn("Failed to set cache for user roles", zap.Int64("uniqueID", uniqueID), zap.Error(err))
		}
	}

	return roles, nil
}

func (s *userService) LoginByEmail(ctx context.Context, data dto.UserLoginByEmailReqDTO) (*dto.UserLoginByEmailResDTO, *exception.Exception) {
	user, err := s.userRepo.FindByEmail(ctx, data.Email)
	if err != nil {
		return nil, exception.ExceptionUserNotFound.AppendDetails(err.Error())
	}

	if !utils.CryptHashCompare(data.Password, user.Salt, user.Password) {
		return nil, exception.ExceptionUserPasswordIncorrect
	}

	token, err := s.jwt.GenerateToken(user.UniqueID)
	if err != nil {
		return nil, exception.ExceptionTokenGenerateFailed.AppendDetails(err.Error())
	}

	return &dto.UserLoginByEmailResDTO{
		Token:     token,
		ExpireAt:  s.jwt.ExpireAt().UnixMilli(),
		RefreshAt: s.jwt.RefreshAt().UnixMilli(),
	}, nil
}
