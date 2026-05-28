package winner

type DrawRequest struct {
	CampaignID string `json:"campaign_id" binding:"required" example:"018f3c4a-2d4c-7b1e-9c1a-2b1c0f2a9f3a"`

	PrizeID string `json:"prize_id" binding:"required" example:"018f3c4a-2d4c-7b1e-9c1a-2b1c0f2a9f3b"`
}

type ListWinnersRequest struct {
	CampaignID string `form:"campaign_id"`
	PrizeID    string `form:"prize_id"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Items []*Winner      `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}