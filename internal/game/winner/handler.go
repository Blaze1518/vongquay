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

// ListWinners godoc
// @Summary      Lấy danh sách người trúng giải (Có phân trang)
// @Description  Lấy danh sách lịch sử trúng thưởng dựa theo chiến dịch hoặc giải thưởng, hỗ trợ phân trang chuẩn hóa.
// @Tags         Winners
// @Accept       json
// @Produce      json
// @Param        campaign_id  query  string  false  "ID của chiến dịch"
// @Param        prize_id     query  string  false  "ID của giải thưởng"
// @Param        page         query  int     false  "Số trang hiện tại (Mặc định: 1)"
// @Param        limit        query  int     false  "Số lượng phần tử mỗi trang (Mặc định: 10)"
// @Success      200          {object}  httperrs.APIError{data=PaginatedResponse} "Thành công"
// @Failure      400          {object}  httperrs.APIError  "Lỗi dữ liệu đầu vào"
// @Failure      500          {object}  httperrs.APIError  "Lỗi hệ thống"
// @Router       /game/winner [get]
func (h *Handler) ListWinners(c *gin.Context) {
	var req ListWinnersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(httperrs.FromGinValidation(err))
		return
	}

	result, err := h.winnerService.ListWinners(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(httperrs.ValidationError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, httperrs.Success(*result))
}