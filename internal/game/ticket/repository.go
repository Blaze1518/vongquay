package ticket

import (
	"context"
	"crypto/rand"
	"math/big"

	baseRepo "github.com/Blaze1518/sinhnhatf168/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(ctx context.Context, ticket *Ticket) error
	CreateInBatchesWithCount(ctx context.Context, tickets []Ticket) (int, error)
	Update(ctx context.Context, ticket *Ticket) error
	CancelOtherTickets(ctx context.Context, campaignID string, username string, winningTicketID string) error
	GetRandomValidTicket(ctx context.Context, campaignID string) (*Ticket, error)
	
	CreateJob(ctx context.Context, job *TicketImportJob) error
	UpdateJobStatus(ctx context.Context, jobID string, status string) error
	UpdateJobProgress(ctx context.Context, jobID string, total, success, failed int) error
	CompleteJob(ctx context.Context, jobID string, total, success, failed int, errorLog string) error
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

func (r *repository) Create(ctx context.Context, ticket *Ticket) error {
	result := r.getDB(ctx).WithContext(ctx).Create(ticket)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository) CreateInBatchesWithCount(ctx context.Context, tickets []Ticket) (int, error) {
	result := r.getDB(ctx).WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns: []clause.Column{
                {Name: "campaign_id"},
                {Name: "ticket_number"},
            },
            DoNothing: true,
        }).
        Create(&tickets) 

	if result.Error != nil {
		return 0, result.Error
	}
	
	return int(result.RowsAffected), nil
}

func (r *repository) Update(ctx context.Context, ticket *Ticket) error {
	return r.getDB(ctx).WithContext(ctx).Save(ticket).Error
}

func (r *repository) CancelOtherTickets(ctx context.Context, campaignID string, username string, winningTicketID string) error {
	return r.getDB(ctx).WithContext(ctx).Model(&Ticket{}).
		Where("campaign_id = ? AND username = ? AND id != ?", campaignID, username, winningTicketID).
		Update("is_canceled", true).Error
}

func (r *repository) CreateJob(ctx context.Context, job *TicketImportJob) error {
	return r.getDB(ctx).WithContext(ctx).Create(job).Error
}

func (r *repository) UpdateJobStatus(ctx context.Context, jobID string, status string) error {
	return r.getDB(ctx).WithContext(ctx).Model(&TicketImportJob{}).
		Where("id = ?", jobID).Update("status", status).Error
}

func (r *repository) UpdateJobProgress(ctx context.Context, jobID string, total, success, failed int) error {
	return r.getDB(ctx).WithContext(ctx).Model(&TicketImportJob{}).Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"total_rows":   total,
			"success_rows": success,
			"failed_rows":  failed,
		}).Error
}

func (r *repository) CompleteJob(ctx context.Context, jobID string, total, success, failed int, errorLog string) error {
	return r.getDB(ctx).WithContext(ctx).Model(&TicketImportJob{}).Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status":       "COMPLETED",
			"total_rows":   total,
			"success_rows": success,
			"failed_rows":  failed,
			"error_log":    errorLog,
		}).Error
}

func (r *repository) GetRandomValidTicket(ctx context.Context, campaignID string) (*Ticket, error) {
	var count int64

	err := r.getDB(ctx).WithContext(ctx).Model(&Ticket{}).
		Where("campaign_id = ? AND is_winner = FALSE AND is_canceled = FALSE", campaignID).
		Count(&count).Error
	
	if err != nil {
		return nil, err
	}
		
	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	bigMax := big.NewInt(count)
	randomBigInt, err := rand.Int(rand.Reader, bigMax)

	if err != nil {
		return nil, err
	}
	offset := int(randomBigInt.Int64())

	var luckyTicket Ticket

	err = r.getDB(ctx).WithContext(ctx).
		Where("campaign_id = ? AND is_winner = FALSE AND is_canceled = FALSE", campaignID).
		Offset(offset).
		Limit(1).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&luckyTicket).Error

	if err != nil {
		return nil, err
	}
	
	return &luckyTicket, nil
}