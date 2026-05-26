package prize

import (
	"net/http"

	httperrs "github.com/Blaze1518/sinhnhatf168/internal/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	prizeService Service
}

func NewHandler(prizeService Service) *Handler {
	return &Handler{prizeService: prizeService}
}

// Create godoc
// @Summary      Tạo mới một giải thưởng (Prize)
// @Description  Tạo mới một giải thưởng thuộc về một chiến dịch cụ thể bằng campaign_id.
// @Tags         Prizes
// @Accept       json
// @Produce      json
// @Param        request body      CreatePrizeRequest  true  "Thông tin giải thưởng cần tạo mới"
// @Success      201     {object}  Prize               "Tạo giải thưởng thành công"
// @Failure      400     {object}  httperrs.APIError   "Dữ liệu gửi lên không hợp lệ hoặc sai định dạng JSON"
// @Failure      500     {object}  httperrs.APIError   "Lỗi hệ thống nội bộ"
// @Router       /game/prize [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreatePrizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperrs.FromGinValidation(err))
		return
	}

	r, err := h.prizeService.Create(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(httperrs.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, httperrs.Success(Prize{
		ID:         r.ID,
		CampaignID: r.CampaignID,
		Name:       r.Name,
		Quantity:   r.Quantity,
		Priority:   r.Priority,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}))
}