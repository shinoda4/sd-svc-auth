package http

import (
	"context"
	stdhttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shinoda4/sd-svc-auth/internal/dto"
)

func (s *Server) HandleRefresh(c *gin.Context) {
	var body dto.RefreshRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	newAccess, accessTTL, err := s.Auth.Refresh(ctx, body.RefreshToken)
	if err != nil {
		c.JSON(stdhttp.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(stdhttp.StatusOK, dto.RefreshResponse{
		AccessToken: newAccess,
		ExpiresIn:   int(accessTTL.Seconds()),
	})
}

func (s *Server) HandleVerifyToken(c *gin.Context) {
	var body dto.VerifyTokenRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	claims, err := s.Auth.ValidateToken(ctx, body.AccessToken)
	if err != nil {
		c.JSON(stdhttp.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(stdhttp.StatusOK, dto.VerifyTokenResponse{
		Token: claims,
	})

}
