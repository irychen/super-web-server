package app

import (
	"context"
	"net/http"
	v1 "super-web-server/internal/api/v1"
	"super-web-server/internal/config"
	"super-web-server/internal/controller"
	"super-web-server/internal/middleware"
	"super-web-server/internal/repo"
	"super-web-server/internal/service"
	"super-web-server/internal/validator"
	"super-web-server/pkg/database"
	"super-web-server/pkg/jwt"
	"super-web-server/pkg/logger"
	"super-web-server/pkg/snowflake"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type App struct {
	config     *config.Config
	engine     *gin.Engine
	server     *http.Server
	db         *database.DB
	redis      *redis.Client
	repo       repo.Repo
	service    service.Service
	controller controller.Controller
	snowflake  *snowflake.Snowflake
	jwt        *jwt.JWT
	roleCheck  *middleware.RoleCheck
}

func NewApp(config *config.Config) (*App, error) {
	app := &App{config: config}
	app.InitEngineAndServer()

	validator.Init()

	if err := app.InitSnowflake(); err != nil {
		return nil, err
	}

	if err := app.InitDatabase(); err != nil {
		return nil, err
	}

	var redisCtx = context.Background()
	redisCtx, cancel := context.WithTimeout(redisCtx, 5*time.Second)
	defer cancel()

	if err := app.InitRedis(redisCtx); err != nil {
		return nil, err
	}

	jwtConfig := app.config.JWT

	app.jwt = jwt.NewJWT(jwt.Config{
		Secret: jwtConfig.Secret,
		Expire: jwtConfig.Expire,
		Issuer: jwtConfig.Issuer,
	})

	app.repo = repo.NewRepo(app.db.DB, logger.GetModuleLogger("repo"))
	app.service = service.NewService(app.repo, logger.GetModuleLogger("service"), app.redis, app.jwt)
	app.roleCheck = middleware.NewRoleCheck(app.service)
	app.controller = controller.NewController(app.service, logger.GetModuleLogger("controller"), app.jwt)

	v1.InitApi(app.engine.Group("api/v1"), app.controller, app.jwt, app.roleCheck)

	return app, nil
}

func (a *App) Run() error {
	logger.InfoF("Starting server on http://localhost:%d", a.config.Server.Port)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
