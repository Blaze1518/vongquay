package ticket

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
)

type Service interface {
	Create(ctx context.Context, req CreateTicketRequest) (*Ticket, error)
	ListTickets(ctx context.Context, req ListTicketsRequest) (*PaginatedResponse, error)
}

type service struct {
	ticketRepo Repository
}

func NewService(ticketRepo Repository) Service {
	return &service{ticketRepo: ticketRepo}
}

func (s *service) Create(ctx context.Context, req CreateTicketRequest) (*Ticket, error) {
	ticketNumber, err := s.generateUniqueTicketNumber()
	if err != nil {
		return nil, fmt.Errorf("không thể sinh mã vé số: %w", err)
	}

	var ticket = &Ticket{
		CampaignID:   req.CampaignID,
		TicketNumber: ticketNumber,
		IsWinner:     false,
	}

	err = s.ticketRepo.Create(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("Có lỗi khi tạo Vé số: %w", err)
	}

	return ticket, nil
}

func (s *service) ListTickets(ctx context.Context, req ListTicketsRequest) (*PaginatedResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	items, totalItems, err := s.ticketRepo.List(ctx, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(totalItems) / req.Limit
	if int(totalItems)%req.Limit != 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Items: items,
		Meta: PaginationMeta{
			CurrentPage: req.Page,
			PerPage:     req.Limit,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

func (s *service) generateUniqueTicketNumber() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 21
	
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	
	return string(result), nil
}