package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleVerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}

	sendEmail := c.DefaultQuery("sendEmail", "true") // 默认 true
	// sendEmail 是字符串，需要转换为 bool
	sendEmailBool := sendEmail == "true"

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := s.Auth.VerifyEmail(ctx, token, sendEmailBool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}
