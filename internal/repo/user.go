package repo

import (
	"context"
	"super-web-server/internal/dto"
	"super-web-server/internal/model"
	"super-web-server/pkg/logger"

	"gorm.io/gorm"
)

type UserRepo interface {
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	Create(ctx context.Context, entity *model.User) error
	Update(ctx context.Context, entity *model.User) error
	SoftDelete(ctx context.Context, id uint64) error
	HardDelete(ctx context.Context, id uint64) error

	FindOne(ctx context.Context, opts ...QueryOption) (*model.User, error)
	FindMany(ctx context.Context, opts ...QueryOption) ([]*model.User, error)
	FindPage(ctx context.Context, pagination dto.Pagination, opts ...QueryOption) ([]*model.User, int64, error)

	UpdateForce(ctx context.Context, entity *model.User) error
	UpdateByMap(ctx context.Context, id uint64, data map[string]any) error

	FindByUniqueID(ctx context.Context, uniqueID int64) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	WithTx(tx *gorm.DB) UserRepo
}

type userRepo struct {
	BaseRepo[model.User]
	db     *gorm.DB
	logger *logger.Logger
}

func NewUserRepo(db *gorm.DB, logger *logger.Logger) UserRepo {
	logger.Info("NewUserRepo initialized successfully")
	return &userRepo{
		BaseRepo: NewBaseRepo[model.User](db, logger),
		db:       db,
		logger:   logger,
	}
}

func (r *userRepo) FindByUniqueID(ctx context.Context, uniqueID int64) (*model.User, error) {
	var opts = []QueryOption{
		Preload("Roles"),
		Where("unique_id = ?", uniqueID),
	}
	return r.BaseRepo.FindOne(ctx, opts...)
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var opts = []QueryOption{
		Preload("Roles"),
		Where("email = ?", email),
	}
	return r.BaseRepo.FindOne(ctx, opts...)
}

func (r *userRepo) WithTx(tx *gorm.DB) UserRepo {
	return &userRepo{
		BaseRepo: r.BaseRepo.WithTx(tx),
		db:       tx,
		logger:   r.logger,
	}
}
