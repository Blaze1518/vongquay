package winner

import (
	"context"

	baseRepo "github.com/Blaze1518/sinhnhatf168/internal/repository"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, winner *Winner) error
	GetWinnersCountByPrize(ctx context.Context, campaignID, prizeID string) (int64, error)
	List(ctx context.Context, campaignID, prizeID string, offset, limit int) ([]*Winner, int64, error)
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

func (r *repository) List(ctx context.Context, campaignID, prizeID string, offset, limit int) ([]*Winner, int64, error) {
	var winners []*Winner
	var total int64

	query := r.getDB(ctx).WithContext(ctx).Model(&Winner{})

	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}
	if prizeID != "" {
		query = query.Where("prize_id = ?", prizeID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&winners).Error

	return winners, total, err
}