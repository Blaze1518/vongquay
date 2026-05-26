package campaign

import "time"

const (
	StatusDraft  = "DRAFT"
	StatusActive = "ACTIVE"
	StatusPaused = "PAUSED"
	StatusEnded  = "ENDED"
)

type Campaign struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
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