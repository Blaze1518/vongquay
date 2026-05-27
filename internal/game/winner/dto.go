package winner

type DrawRequest struct {
	CampaignID string `json:"campaign_id" binding:"required" example:"018f3c4a-2d4c-7b1e-9c1a-2b1c0f2a9f3a"`

	PrizeID string `json:"prize_id" binding:"required" example:"018f3c4a-2d4c-7b1e-9c1a-2b1c0f2a9f3b"`
}