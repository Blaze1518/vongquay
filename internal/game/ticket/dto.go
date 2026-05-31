package ticket

import (
	"errors"
	"mime/multipart"
	"path/filepath"
)

type CreateTicketRequest struct {
	CampaignID string `json:"campaign_id" binding:"required" example:"018f3c4a-2d4c-7b1e-9c1a-2b1c0f2a9f3a"`
}

type UpdateTicketRequest struct {
	IsWinner *bool `json:"is_winner" binding:"omitempty" example:"true"`
}

type ListTicketsRequest struct {
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
	Items []*Ticket      `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type ImportExcelRequest struct {
	CampaignID string                `form:"campaign_id" binding:"required" example:"018f3c4a-2d4c-7b1e-9c1a-2b1c0f2a9f3b"`
	File       *multipart.FileHeader `form:"file" binding:"required"`
}

func (r *ImportExcelRequest) Validate() error {
	ext := filepath.Ext(r.File.Filename)
	if ext != ".xlsx" {
		return errors.New("hệ thống chỉ chấp nhận định dạng file .xlsx")
	}

	const maxFileSize = 20 * 1024 * 1024
	if r.File.Size > maxFileSize {
		return errors.New("dung lượng file vượt quá giới hạn cho phép (Tối đa 20MB)")
	}

	return nil
}