package winner

import "time"

type Winner struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CampaignID uint      `gorm:"index;not null" json:"campaign_id"`
	PrizeID    uint      `gorm:"index;not null" json:"prize_id"`
	TicketID   uint      `gorm:"index;not null" json:"ticket_id"`
	DrawOrder  int       `json:"draw_order"`

	CreatedAt  time.Time `json:"created_at"`
}

func (Winner) TableName() string {
	return "winners"
}