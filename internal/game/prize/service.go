package prize

import (
	"context"
	"fmt"
)

type Service interface {
	Create(ctx context.Context, req CreatePrizeRequest) (*Prize, error)
}

type service struct {
	prizeRepo Repository
}

func NewService(prizeRepo Repository) Service {
	return &service{prizeRepo: prizeRepo}
}

func (s *service) Create(ctx context.Context, req CreatePrizeRequest) (*Prize, error) {
	var prize = &Prize{
		CampaignID: req.CampaignID,
		Name:       req.Name,
		Quantity:   req.Quantity,
		Priority:   req.Priority,
	}

	if prize.Quantity <= 0 {
		prize.Quantity = 1
	}

	err := s.prizeRepo.Create(ctx, prize)
	if err != nil {
		return nil, fmt.Errorf("Có lỗi khi tạo Giải thưởng: %w", err)
	}

	return prize, nil
}