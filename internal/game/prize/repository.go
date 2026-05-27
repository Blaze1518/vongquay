package prize

import (
	"context"

	baseRepo "github.com/Blaze1518/sinhnhatf168/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(ctx context.Context, prize *Prize) error
	Update(ctx context.Context, prize *Prize) error
	GetPrizeForUpdate(ctx context.Context, id string) (*Prize, error)
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

func (r *repository) Create(ctx context.Context, prize *Prize) error {
	result := r.getDB(ctx).WithContext(ctx).Create(prize)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository) GetPrizeForUpdate(ctx context.Context, id string) (*Prize, error) {
	var p Prize
	err := r.getDB(ctx).WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&p).Error

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) Update(ctx context.Context, prize *Prize) error {
	return r.getDB(ctx).WithContext(ctx).Save(prize).Error
}