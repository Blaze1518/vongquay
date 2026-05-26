package campaign

import (
	"errors"
	"net/http"

	httperrs "github.com/Blaze1518/sinhnhatf168/internal/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	campaignService Service
}

func NewHandler(campaignService Service) *Handler {
	return &Handler{campaignService: campaignService}
}

// Create godoc
// @Summary      Tạo mới một chiến dịch (Campaign)
// @Description  Tạo mới một chiến dịch với tên và mã code riêng biệt. Hệ thống tự động chặn trùng mã code.
// @Tags         Campaigns
// @Accept       json
// @Produce      json
// @Param        request body      CreateCampaignRequest  true  "Thông tin chiến dịch cần tạo mới"
// @Success      201     {object}  Campaign               "Tạo chiến dịch thành công"
// @Failure      400     {object}  httperrs.APIError      "Dữ liệu gửi lên không hợp lệ hoặc sai định dạng JSON"
// @Failure      409     {object}  httperrs.APIError      "Mã định danh chiến dịch (Code) đã tồn tại"
// @Failure      500     {object}  httperrs.APIError      "Lỗi hệ thống nội bộ"
// @Router       /game/campaign [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperrs.FromGinValidation(err))
		return
	}

	r, err := h.campaignService.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrCampaignExists) {
			_ = c.Error(httperrs.Conflict("Mã chiến dịch đã tồn tại"))
			return
		}
		_ = c.Error(httperrs.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, httperrs.Success(Campaign{
		ID: r.ID,
		Name: r.Name,
		Code: r.Code,
		Status: r.Status,
		StartedAt: r.StartedAt,
		EndedAt: r.EndedAt,
	}))
}