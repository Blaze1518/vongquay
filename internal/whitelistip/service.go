package whitelistip

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrWhitelistIPExists = errors.New("IP đã tồn tại")
	ErrWhitelistIPNotFound = errors.New("IP không tồn tại")
)

type Service interface {
	CreateWhitelistIP(ctx context.Context, req CreateWhitelistIPRequest) (*WhitelistIP, error)
	IsWhitelistedIP(ctx context.Context, ip string) (bool, error)
}

type service struct {
	whitelistIPRepo Repository
}

func NewService(whitelistIPRepo Repository) Service {
	return &service{whitelistIPRepo: whitelistIPRepo}
}

func (s *service) CreateWhitelistIP(ctx context.Context, req CreateWhitelistIPRequest) (*WhitelistIP, error) {
	var whitelistIP = &WhitelistIP{
		IPAddress: req.IPAddress,
		Description: req.Description,
		IsActive: *req.IsActive,
	}

	err := s.whitelistIPRepo.Create(ctx, whitelistIP)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrWhitelistIPExists
		}
		return nil, fmt.Errorf("Có lỗi khi tạo IP: %w", err)
	}

	return whitelistIP, nil
}

func (s *service) IsWhitelistedIP(ctx context.Context, ip string) (bool, error) {
	whitelistIP, err := s.whitelistIPRepo.FindByIPAddress(ctx, ip)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrWhitelistIPNotFound
		}
		return false, fmt.Errorf("Có lỗi khi tìm kiếm IP: %w", err)
	}
	return whitelistIP.IsActive, nil
}