package middleware

import "github.com/gin-gonic/gin"

func ResolveClientIP(c *gin.Context) string {
    ip := c.ClientIP()
    if ip == "" {
        return "unknown"
    }
    return ip
}