package prize

type CreatePrizeRequest struct {
	CampaignID uint   `json:"campaign_id" binding:"required" example:"1"`
	Name       string `json:"name" binding:"required,max=255" example:"Xe máy Honda Vision 2026"`
	Quantity   int    `json:"quantity" binding:"required,min=1" example:"10"`
	Priority   int    `json:"priority" example:"1"`
}

type UpdatePrizeRequest struct {
	Name     string `json:"name" binding:"omitempty,max=255" example:"Xe máy Honda Vision 2026 (Cập nhật)"`
	Quantity *int   `json:"quantity" binding:"omitempty,min=0" example:"15"`
	Priority *int   `json:"priority" binding:"omitempty" example:"2"`
}