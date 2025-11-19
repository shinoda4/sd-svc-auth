package http

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if len(h) < 7 || h[:7] != "Bearer " {
			c.AbortWithStatusJSON(stdhttp.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}
		tokenStr := h[7:]
		claims, err := s.Auth.ValidateToken(c.Request.Context(), tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(stdhttp.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
