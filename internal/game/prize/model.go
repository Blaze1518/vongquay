package prize

import "time"

type Prize struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CampaignID uint      `gorm:"index;not null" json:"campaign_id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Quantity   int       `gorm:"default:1" json:"quantity"`
	Priority   int       `json:"priority"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Prize) TableName() string {
	return "prizes"
}