package whitelistip

import (
	"context"

	baseRepo "github.com/Blaze1518/sinhnhatf168/internal/repository"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, whitelistIP *WhitelistIP) error
	FindByIPAddress(ctx context.Context, ipAddress string) (*WhitelistIP, error)
	Transaction(ctx context.Context, fn func(context.Context) error) error
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

func (r *repository) Create(ctx context.Context, whitelistIP *WhitelistIP) error {
	result := r.getDB(ctx).WithContext(ctx).Create(whitelistIP)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository) FindByIPAddress(ctx context.Context, ipAddress string) (*WhitelistIP, error) {
	var whitelistIP WhitelistIP
	result := r.getDB(ctx).WithContext(ctx).Where("ip_address = ?", ipAddress).First(&whitelistIP)
	if result.Error != nil {
		return nil, result.Error
	}
	return &whitelistIP, nil
}

func (r *repository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return baseRepo.RunInTransaction(r.db, ctx, fn)
}