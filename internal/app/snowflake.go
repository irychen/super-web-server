package app

import (
	"fmt"
	"super-web-server/pkg/logger"
	"super-web-server/pkg/snowflake"
)

func (a *App) InitSnowflake() error {
	snowflake, err := snowflake.NewSnowflake(a.config.Server.SnowflakeNode)
	if err != nil {
		return fmt.Errorf("init snowflake failed %w", err)
	}
	a.snowflake = snowflake

	logger.Info("snowflake initialized successfully")
	return nil
}
