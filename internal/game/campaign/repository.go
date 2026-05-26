package campaign

import (
	"context"

	baseRepo "github.com/Blaze1518/sinhnhatf168/internal/repository"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, campaign *Campaign) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) getDB(ctx context.Context) *gorm.DB {
	return baseRepo.GetDB(ctx, r.db)
}

func (r *repository) Create(ctx context.Context, campaign *Campaign) error {
	result := r.getDB(ctx).WithContext(ctx).Create(campaign)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

