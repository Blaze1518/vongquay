package whitelistip

import (
	"errors"

	httperrs "github.com/Blaze1518/sinhnhatf168/internal/errors"
	"github.com/gin-gonic/gin"
)

func WhitelistIPMiddleware(
	ip func(c *gin.Context) string,
	whitelistIPService Service,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := ip(c)
		isWhitelisted, err := whitelistIPService.IsWhitelistedIP(c.Request.Context(), ip)
		if err != nil {
			if errors.Is(err, ErrWhitelistIPNotFound) {
				_ = c.Error(httperrs.Forbidden("IP không nằm trong danh sách truy cập"))
				c.Abort()
				return
			}
			_ = c.Error(httperrs.InternalServerError(err))
			c.Abort()
			return
		}
		if !isWhitelisted {
			_ = c.Error(httperrs.Forbidden("IP chưa được kích hoạt, vui lòng liên hệ admin để được kích hoạt"))
			c.Abort()
			return
		}
		c.Next()
	}
}