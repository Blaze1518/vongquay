package whitelistip

import (
	"errors"
	"net/http"

	httperrs "github.com/Blaze1518/sinhnhatf168/internal/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	whitelistIPService Service
}

func NewHandler(whitelistIPService Service) *Handler {
	return &Handler{whitelistIPService: whitelistIPService}
}

// CreateWhitelistIP godoc
// @Id           createWhitelistIP
// @Summary      Create whitelist IP
// @Description  Creates a unique IP or CIDR row in table whitelist_ips. Stored as PostgreSQL cidr: use a network prefix (e.g. 203.0.113.10/32 for a single IPv4 host, 2001:db8::/32 for an IPv6 range). Duplicate ip_address returns 409. is_active is stored as-is; access control must ignore rows with is_active=false.
// @Description  On success the JSON envelope uses success/data; data is a WhitelistIP object (id, ip_address, description, is_active, created_at, updated_at).
// @Tags         whitelist-ips
// @Accept       json
// @Produce      json
// @Param        body body CreateWhitelistIPRequest true "Create whitelist IP payload"
// @Success      201 {object} httperrs.Response "Created; data field contains WhitelistIP"
// @Failure      400 {object} httperrs.Response "Validation error (invalid JSON or binding rules)"
// @Failure      409 {object} httperrs.Response "Conflict — IP/CIDR already exists"
// @Failure      500 {object} httperrs.Response "Unexpected server error"
// @Router       /whitelist-ips/ [post]
func (h *Handler) CreateWhitelistIP(c *gin.Context) {
	var req CreateWhitelistIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperrs.FromGinValidation(err))
		return
	}

	whitelistIP, err := h.whitelistIPService.CreateWhitelistIP(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrWhitelistIPExists) {
			_ = c.Error(httperrs.Conflict("IP đã tồn tại"))
			return
		}
		_ = c.Error(httperrs.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, httperrs.Success(WhitelistIP{
		ID: whitelistIP.ID,
		IPAddress: whitelistIP.IPAddress,
		Description: whitelistIP.Description,
		IsActive: whitelistIP.IsActive,
		CreatedAt: whitelistIP.CreatedAt,
		UpdatedAt: whitelistIP.UpdatedAt,
	}))
}