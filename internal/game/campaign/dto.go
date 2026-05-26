package campaign

import "time"

type CreateCampaignRequest struct {
	Name      string     `json:"name" binding:"required,max=255" example:"Chương trình Khuyến mãi Sinh nhật F168"`
	Code      string     `json:"code" binding:"required,max=100" example:"SINHNHAT_F168_2026"`
	Status    string     `json:"status" binding:"omitempty,oneof=DRAFT ACTIVE PAUSED ENDED" example:"ACTIVE"`
	StartedAt *time.Time `json:"started_at,omitempty" example:"2026-05-26T00:00:00Z"`
	EndedAt   *time.Time `json:"ended_at,omitempty" example:"2026-06-26T23:59:59Z"`
}

type UpdateCampaignRequest struct {
	Name      string `json:"name" binding:"omitempty,max=255" example:"Chương trình Khuyến mãi Sinh nhật F168 (Cập nhật)"`
	Code      string `json:"code" binding:"omitempty,max=100" example:"SINHNHAT_F168_2026_V2"`
	Status    string `json:"status" binding:"omitempty,oneof=DRAFT ACTIVE PAUSED ENDED" example:"PAUSED"`
	StartedAt *time.Time `json:"started_at,omitempty" example:"2026-05-27T00:00:00Z"`
	EndedAt   *time.Time `json:"ended_at,omitempty" example:"2026-06-27T23:59:59Z"`
}