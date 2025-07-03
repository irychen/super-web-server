package repo

import (
	"super-web-server/pkg/logger"

	"gorm.io/gorm"
)

type Repo interface {
	User() UserRepo
}

type repo struct {
	userRepo UserRepo
	logger   *logger.Logger
}

func NewRepo(db *gorm.DB, logger *logger.Logger) Repo {
	logger.Info("NewRepo initialized successfully")
	return &repo{
		userRepo: NewUserRepo(db, logger),
		logger:   logger,
	}
}

func (r *repo) User() UserRepo {
	return r.userRepo
}
