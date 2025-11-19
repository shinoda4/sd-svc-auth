package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shinoda4/sd-svc-auth/internal/dto"
)

func (s *Server) HandlePasswordReset(c *gin.Context) {
	// sendEmail := c.DefaultQuery("sendEmail", "true")
	// sendEmailBool := sendEmail == "true"

	var body dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "empty request body",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := s.Auth.PasswordReset(ctx, body.Email, body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

func (s *Server) HandlePasswordResetConfirm(c *gin.Context) {

	token := c.DefaultQuery("token", "")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token not provided",
		})
		return
	}

	var body dto.ResetPasswordConfirmRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "empty request body",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if body.NewPassword != body.NewPasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password not confirm",
		})
		return
	}

	err := s.Auth.PasswordResetConfirm(ctx, token, body.NewPasswordConfirm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
