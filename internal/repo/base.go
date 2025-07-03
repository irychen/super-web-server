package repo

import (
	"context"
	"super-web-server/internal/dto"
	"super-web-server/pkg/logger"

	"gorm.io/gorm"
)

type BaseRepo[T any] interface {
	// basic
	FindByID(ctx context.Context, id uint64) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	SoftDelete(ctx context.Context, id uint64) error
	HardDelete(ctx context.Context, id uint64) error

	// find
	FindOne(ctx context.Context, opts ...QueryOption) (*T, error)
	FindMany(ctx context.Context, opts ...QueryOption) ([]*T, error)
	FindPage(ctx context.Context, pagination dto.Pagination, opts ...QueryOption) ([]*T, int64, error)

	// special update
	UpdateForce(ctx context.Context, entity *T) error
	UpdateByMap(ctx context.Context, id uint64, data map[string]any) error

	WithTx(tx *gorm.DB) BaseRepo[T]
}

type baseRepo[T any] struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewBaseRepo[T any](db *gorm.DB, logger *logger.Logger) BaseRepo[T] {
	return &baseRepo[T]{db: db, logger: logger}
}

func (r *baseRepo[T]) FindByID(ctx context.Context, id uint64) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepo[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepo[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *baseRepo[T]) SoftDelete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(new(T), id).Error
}

func (r *baseRepo[T]) HardDelete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Unscoped().Delete(new(T), id).Error
}

func (r *baseRepo[T]) FindOne(ctx context.Context, opts ...QueryOption) (*T, error) {
	var entity T
	db := ApplyQueryOptions(r.db.WithContext(ctx), opts...)
	if err := db.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepo[T]) FindMany(ctx context.Context, opts ...QueryOption) ([]*T, error) {
	var entities = make([]*T, 0)
	db := ApplyQueryOptions(r.db.WithContext(ctx), opts...)
	if err := db.Find(&entities).Error; err != nil {
		return entities, err
	}
	return entities, nil
}

func (r *baseRepo[T]) FindPage(ctx context.Context, pagination dto.Pagination, opts ...QueryOption) ([]*T, int64, error) {
	var entities []*T
	var total int64
	offset, limit := pagination.Offset(), pagination.Limit()

	db := ApplyQueryOptions(r.db.WithContext(ctx), opts...)

	var entity T
	err := db.Model(&entity).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset(offset).Limit(limit).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *baseRepo[T]) UpdateForce(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Model(entity).Select("*").Updates(entity).Error
}

func (r *baseRepo[T]) UpdateByMap(ctx context.Context, id uint64, data map[string]any) error {
	return r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(data).Error
}

func (r *baseRepo[T]) WithTx(tx *gorm.DB) BaseRepo[T] {
	return &baseRepo[T]{db: tx}
}
