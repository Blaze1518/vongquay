package winner

import (
	"context"

	baseRepo "github.com/Blaze1518/sinhnhatf168/internal/repository"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, winner *Winner) error
	GetWinnersCountByPrize(ctx context.Context, campaignID, prizeID string) (int64, error)
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

func (r *repository) Create(ctx context.Context, winner *Winner) error {
	return r.getDB(ctx).WithContext(ctx).Create(winner).Error
}

func (r *repository) GetWinnersCountByPrize(ctx context.Context, campaignID, prizeID string) (int64, error) {
	var count int64
	err := r.getDB(ctx).WithContext(ctx).Model(&Winner{}).
		Where("campaign_id = ? AND prize_id = ?", campaignID, prizeID).
		Count(&count).Error
	return count, err
}