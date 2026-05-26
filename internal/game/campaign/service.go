package campaign

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrCampaignExists = errors.New("Chiến dịch đã tồn tại")
	ErrCampaignNotFound = errors.New("Chiến dịch không tồn tại")
)

type Service interface {
	Create(ctx context.Context, req CreateCampaignRequest) (*Campaign, error)
}

type service struct {
	campaignRepo Repository
}

func NewService(campaignRepo Repository) Service {
	return &service{campaignRepo: campaignRepo}
}

func (s *service) Create(ctx context.Context, req CreateCampaignRequest) (*Campaign, error) {
	var campaign = &Campaign{
		Name: req.Name,
		Code: req.Code,
		Status: req.Status,
		StartedAt: req.StartedAt,
		EndedAt: req.EndedAt,
	}

	err := s.campaignRepo.Create(ctx, campaign)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrCampaignExists
		}
		return nil, fmt.Errorf("Có lỗi khi tạo Campaign: %w", err)
	}

	return campaign, nil
}