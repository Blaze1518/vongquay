package whitelistip

import "time"

// WhitelistIP is the persisted row and the shape of `data` on successful create.
type WhitelistIP struct {
	ID          uint      `gorm:"primaryKey" json:"id" example:"1"`
	IPAddress   string    `gorm:"type:cidr;not null;uniqueIndex" json:"ip_address" example:"203.0.113.10/32"`
	Description string    `gorm:"type:varchar(255)" json:"description" example:"HQ VPN egress - INC-1234"`
	IsActive    bool      `gorm:"default:true" json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2026-04-26T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-04-26T12:00:00Z"`
}

func (WhitelistIP) TableName() string {
	return "whitelist_ips"
}