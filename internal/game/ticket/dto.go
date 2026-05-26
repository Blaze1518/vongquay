package ticket

import (
	"errors"
	"mime/multipart"
	"path/filepath"
)

type CreateTicketRequest struct {
	CampaignID uint `json:"campaign_id" binding:"required" example:"1"`
}

type UpdateTicketRequest struct {
	IsWinner *bool `json:"is_winner" binding:"omitempty" example:"true"`
}

type ImportExcelRequest struct {
	CampaignID uint                  `form:"campaign_id" binding:"required" example:"1"`
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