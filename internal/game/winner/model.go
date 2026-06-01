package winner

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Winner struct {
	ID        string     `gorm:"type:uuid;primaryKey" json:"id"`
	CampaignID   string    `gorm:"type:uuid;index;not null" json:"campaign_id"`
	PrizeID      string    `gorm:"type:uuid;index;not null" json:"prize_id"`
	PrizeName    string    `gorm:"size:255;not null" json:"prize_name"`
	TicketID     string    `gorm:"type:uuid;index;not null" json:"ticket_id"`
	TicketNumber string    `gorm:"size:21;not null" json:"ticket_number"`
	Username     string    `gorm:"size:100;not null" json:"username"`
	DrawOrder    int       `json:"draw_order"`

	CreatedAt  time.Time `json:"created_at"`
}

func (Winner) TableName() string {
	return "winners"
}

func (c *Winner) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	c.ID = id.String()
	return nil
}