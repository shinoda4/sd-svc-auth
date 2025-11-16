package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RefreshBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *Server) HandleRefresh(c *gin.Context) {
	var body RefreshBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	newAccess, accessTTL, err := s.Auth.Refresh(ctx, body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccess,
		"expires_in":   int(accessTTL.Seconds()),
	})
}

func (s *Server) HandleVerifyToken(c *gin.Context) {
	var body struct {
		AccessToken string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	claims, err := s.Auth.ValidateToken(ctx, body.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": claims,
	})

}
