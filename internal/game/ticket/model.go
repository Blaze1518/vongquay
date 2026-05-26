package ticket

import "time"

type Ticket struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CampaignID   uint      `gorm:"index;not null" json:"campaign_id"`
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