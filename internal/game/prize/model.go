package prize

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Prize struct {
	ID        string     `gorm:"type:uuid;primaryKey" json:"id"`
	CampaignID string    `gorm:"type:uuid;index;not null" json:"campaign_id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Quantity   int       `gorm:"default:1" json:"quantity"`
	Priority   int       `json:"priority"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Prize) TableName() string {
	return "prizes"
}

func (c *Prize) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	c.ID = id.String()
	return nil
}