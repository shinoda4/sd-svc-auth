package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ResetPasswordBody struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (s *Server) HandlePasswordReset(c *gin.Context) {
	// sendEmail := c.DefaultQuery("sendEmail", "true")
	// sendEmailBool := sendEmail == "true"

	var body ResetPasswordBody
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

type ResetPasswordConfirmBody struct{
	NewPassword string `json:"new_password" binding:"required"`
	NewPasswordConfirm string `json:"new_password_confirm" binding:"required"`
}
func (s *Server) HandlePasswordResetConfirm(c *gin.Context) {

	token := c.DefaultQuery("token", "")

	if token == ""{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token not provided",
		})
		return
	}

	var body ResetPasswordConfirmBody
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

	if body.NewPassword != body.NewPasswordConfirm{
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
