package ticket

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID        string     `gorm:"type:uuid;primaryKey" json:"id"`
	CampaignID   string    `gorm:"type:uuid;index;not null" json:"campaign_id"`
	TicketNumber string    `gorm:"size:21;not null" json:"ticket_number"`
	Username     string    `gorm:"size:100;not null" json:"username"`
	IsWinner     bool      `gorm:"default:false" json:"is_winner"`
	IsCanceled   bool      `gorm:"default:false" json:"is_canceled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Ticket) TableName() string {
	return "tickets"
}

func (c *Ticket) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	c.ID = id.String()
	return nil
}
