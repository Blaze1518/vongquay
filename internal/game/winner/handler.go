package winner

import (
	"net/http"

	httperrs "github.com/Blaze1518/sinhnhatf168/internal/errors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	winnerService Service
}

func NewHandler(winnerService Service) *Handler {
	return &Handler{winnerService: winnerService}
}

// Draw godoc
// @Summary      Kích hoạt quay số xác định người trúng giải (Draw)
// @Description  Thực hiện bốc thăm ngẫu nhiên bảo mật dưới dạng Transaction để tìm vé trúng thưởng, tự động trừ kho giải và hủy các vé còn lại của user trúng giải.
// @Tags         Winners
// @Accept       json
// @Produce      json
// @Param        request body      DrawRequest        true  "Thông tin chiến dịch và giải thưởng cần quay"
// @Success      200     {object}  Winner             "Quay số thành công và đã xác định người trúng giải"
// @Failure      400     {object}  httperrs.APIError  "Lỗi nghiệp vụ (Hết giải thưởng, Hết vé hợp lệ hoặc sai dữ liệu đầu vào)"
// @Failure      500     {object}  httperrs.APIError  "Lỗi hệ thống nội bộ"
// @Router       /game/winner/draw [post]
func (h *Handler) Draw(c *gin.Context) {
	var req DrawRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperrs.FromGinValidation(err))
		return
	}

	result, err := h.winnerService.ExecuteDraw(c.Request.Context(), req.CampaignID, req.PrizeID)
	if err != nil {
		_ = c.Error(httperrs.ValidationError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, httperrs.Success(*result))
}