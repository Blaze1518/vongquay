package campaign

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusDraft  = "DRAFT"
	StatusActive = "ACTIVE"
	StatusPaused = "PAUSED"
	StatusEnded  = "ENDED"
)

type Campaign struct {
	ID        string     `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	Code      string     `gorm:"size:100;uniqueIndex" json:"code"`
	Status    string     `gorm:"size:30;default:ACTIVE" json:"status"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Campaign) TableName() string {
	return "campaigns"
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	c.ID = id.String()
	return nil
}