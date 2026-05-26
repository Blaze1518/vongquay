package winner

type DrawRequest struct {
	CampaignID uint `json:"campaign_id" binding:"required" example:"1"`

	PrizeID uint `json:"prize_id" binding:"required" example:"2"`
}