package app

import (
	"context"
	"super-web-server/pkg/logger"
	"super-web-server/pkg/redis"
)

func (a *App) InitRedis(ctx context.Context) error {
	redisClient, err := redis.NewRedis(
		&redis.Config{
			Host:     a.config.Redis.Host,
			Port:     a.config.Redis.Port,
			Password: a.config.Redis.Password,
			DB:       a.config.Redis.DB,
		},
		ctx,
	)
	if err != nil {
		return err
	}
	a.redis = redisClient
	logger.Info("redis initialized successfully")
	return nil
}
