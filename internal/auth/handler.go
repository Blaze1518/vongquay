package auth

import (
	"net/http"

	httperrs "github.com/Blaze1518/sinhnhatf168/internal/errors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService Service
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Login(c *gin.Context) {
	ok := "OK"
	c.JSON(http.StatusOK, httperrs.Success(AuthResponse{
		ok: ok,
	}))
}