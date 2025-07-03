package controller

import (
	"super-web-server/internal/service"
	"super-web-server/pkg/jwt"
	"super-web-server/pkg/logger"
)

type Controller interface {
	Hello() HelloController
	User() UserController
}

type controller struct {
	helloController HelloController
	userController  UserController
	logger          *logger.Logger
	jwt             *jwt.JWT
}

func NewController(service service.Service, logger *logger.Logger, jwt *jwt.JWT) Controller {
	logger.Info("NewController initialized successfully")
	return &controller{
		helloController: NewHelloController(logger),
		userController:  NewUserController(service.User(), logger),
		logger:          logger,
		jwt:             jwt,
	}
}

func (c *controller) Hello() HelloController {
	return c.helloController
}

func (c *controller) User() UserController {
	return c.userController
}
