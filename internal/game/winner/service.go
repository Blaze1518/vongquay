package winner

import (
	"context"
	"errors"
	"fmt"

	"github.com/Blaze1518/sinhnhatf168/internal/game/prize"
	"github.com/Blaze1518/sinhnhatf168/internal/game/ticket"
	baseRepo "github.com/Blaze1518/sinhnhatf168/internal/repository"
	"gorm.io/gorm"
)

type Service interface {
	ExecuteDraw(ctx context.Context, campaignID uint, prizeID uint) (*Winner, error)
}

type service struct {
	winnerRepo Repository
	ticketRepo ticket.Repository
	prizeRepo  prize.Repository
	db         *gorm.DB
}

func NewService(wr Repository, tr ticket.Repository, pr prize.Repository, db *gorm.DB) Service {
	return &service{
		winnerRepo: wr,
		ticketRepo: tr,
		prizeRepo:  pr,
		db:         db,
	}
}

func (s *service) ExecuteDraw(ctx context.Context, campaignID uint, prizeID uint) (*Winner, error) {
	var winnerResult *Winner
	err := baseRepo.RunInTransaction(s.db, ctx, func(txCtx context.Context) error {
		
		p, err := s.prizeRepo.GetPrizeForUpdate(txCtx, prizeID)
		if err != nil {
			return fmt.Errorf("không tìm thấy giải thưởng: %w", err)
		}
		if p.Quantity <= 0 {
			return errors.New("giải thưởng này đã được phát hết")
		}

		luckyTicket, err := s.ticketRepo.GetRandomValidTicket(ctx, campaignID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("không còn vé số nào hợp lệ để quay thưởng")
			}
			return fmt.Errorf("lỗi khi quay số tìm vé: %w", err)
		}

		luckyTicket.IsWinner = true
		if err := s.ticketRepo.Update(txCtx, luckyTicket); err != nil {
			return fmt.Errorf("thất bại khi cập nhật vé trúng: %w", err)
		}
		
		if err := s.ticketRepo.CancelOtherTickets(txCtx, campaignID, luckyTicket.Username, luckyTicket.ID); err != nil {
			return fmt.Errorf("thất bại khi hủy các vé còn lại của user: %w", err)
		}

		p.Quantity = p.Quantity - 1
		if err := s.prizeRepo.Update(ctx, p); err != nil {
			return fmt.Errorf("thất bại khi trừ số lượng giải thưởng: %w", err)
		}

		existingWinners, err := s.winnerRepo.GetWinnersCountByPrize(txCtx, campaignID, prizeID)
		
		if err != nil {
			return fmt.Errorf("thất bại khi tính toán số thứ tự lượt quay: %w", err)
		}

		winnerResult = &Winner{
			CampaignID:   campaignID,
			PrizeID:      prizeID,
			TicketID:     luckyTicket.ID,
			DrawOrder:    int(existingWinners) + 1,
		}
		if err := s.winnerRepo.Create(txCtx, winnerResult); err != nil {
			return fmt.Errorf("thất bại khi ghi nhận người trúng giải: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return winnerResult, nil
}