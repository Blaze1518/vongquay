package ticket

import (
	"net/http"

	httperrs "github.com/Blaze1518/sinhnhatf168/internal/errors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ticketService Service
	importService ImportService
}

func NewHandler(ticketService Service, importService ImportService) *Handler {
	return &Handler{ticketService: ticketService, importService: importService}
}

// Create godoc
// @Summary      Phát hành một vé số mới (Ticket)
// @Description  Tạo mới một lượt chơi/vé số thuộc về một chiến dịch. Mã vé dài 21 ký tự sẽ tự động được sinh ngẫu nhiên bảo mật ở hệ thống.
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        request body      CreateTicketRequest  true  "Thông tin chiến dịch để cấp vé"
// @Success      201     {object}  Ticket               "Phát hành vé thành công"
// @Failure      400     {object}  httperrs.APIError    "Dữ liệu gửi lên không hợp lệ hoặc sai định dạng JSON"
// @Failure      500     {object}  httperrs.APIError    "Lỗi hệ thống nội bộ"
// @Router       /game/ticket [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperrs.FromGinValidation(err))
		return
	}

	r, err := h.ticketService.Create(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(httperrs.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, httperrs.Success(Ticket{
		ID:           r.ID,
		CampaignID:   r.CampaignID,
		TicketNumber: r.TicketNumber,
		IsWinner:     r.IsWinner,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}))
}

// ImportExcel godoc
// @Summary      Upload file Excel để import lượng lớn vé số
// @Description  Nhận file .xlsx xử lý bất đồng bộ dưới nền.
// @Tags         Tickets
// @Accept       multipart/form-data
// @Produce      json
// @Param        campaign_id  formData  string   true  "ID của chiến dịch (UUID v7)"
// @Param        file         formData  file  true  "File Excel chứa danh sách mã vé (.xlsx)"
// @Success      202          {object}  TicketImportJob "Đã tiếp nhận file và đang xử lý ngầm"
// @Router       /game/ticket/import [post]
func (h *Handler) ImportExcel(c *gin.Context) {
	var req ImportExcelRequest
	
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(httperrs.FromGinValidation(err))
		return
	}

	if err := req.Validate(); err != nil {
		_ = c.Error(httperrs.ValidationError(err.Error()))
		return
	}

	job, err := h.importService.Import(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(httperrs.InternalServerError(err))
		return
	}

	c.JSON(http.StatusAccepted, httperrs.Success(job))
}

// ListTickets godoc
// @Summary      Lấy danh sách vé số (Phân trang)
// @Description  Trả về danh sách các vé số đã được phát hành trong hệ thống, sắp xếp theo thời gian tạo giảm dần (mới nhất lên đầu). Hỗ trợ phân trang qua query parameters.
// @Tags         Tickets
// @Produce      json
// @Param        request      query     ListTicketsRequest  true  "Tham số phân trang và bộ lọc"
// @Success      200          {object}  PaginatedResponse         "Lấy danh sách vé số thành công"
// @Failure      400          {object}  httperrs.APIError         "Tham số gửi lên không hợp lệ"
// @Failure      500          {object}  httperrs.APIError         "Lỗi hệ thống nội bộ"
// @Router       /game/ticket [get]
func (h *Handler) ListTickets(c *gin.Context) {
    var req ListTicketsRequest

    if err := c.ShouldBindQuery(&req); err != nil {
        _ = c.Error(httperrs.FromGinValidation(err))
        return
    }

    result, err := h.ticketService.ListTickets(c.Request.Context(), req)
    if err != nil {
        _ = c.Error(httperrs.ValidationError(err.Error()))
        return 
    }

    c.JSON(http.StatusOK, httperrs.Success(*result))
}