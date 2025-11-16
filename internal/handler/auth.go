package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shinoda4/sd-svc-auth/internal/repo"
)

type RegisterBody struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
type RegisterResp struct {
	Message     string `json:"message"`
	VerifyToken string `json:"verifyToken"`
}

func (s *Server) HandleRegister(c *gin.Context) {

	sendEmail := c.DefaultQuery("sendEmail", "true") // 默认 true
	// sendEmail 是字符串，需要转换为 bool
	sendEmailBool := sendEmail == "true"

	var body RegisterBody
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

	host := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	verifyLink := fmt.Sprintf("%s://%s/api/v1/verify", scheme, host)

	_, verifyToken, err := s.Auth.Register(ctx, body.Email, body.Username, body.Password, sendEmailBool, verifyLink)
	if err != nil {
		var e *repo.ErrUserExists
		if errors.As(err, &e) {
			c.JSON(http.StatusConflict, gin.H{"error": "register email exists", "details": e.Email})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed", "details": err.Error()})
		}
		return
	}

	resp := RegisterResp{
		Message:     "registered",
		VerifyToken: verifyToken,
	}
	c.JSON(http.StatusCreated, resp)
	return

}

type LoginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}

func (s *Server) HandleLogin(c *gin.Context) {
	var body LoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty request body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accessToken, refreshToken, accessTTL, refreshTTL, err := s.Auth.Login(ctx, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp := LoginResp{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        int(accessTTL.Seconds()),
		RefreshExpiresIn: int(refreshTTL.Seconds()),
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) HandleLogout(c *gin.Context) {
	h := c.GetHeader("Authorization")
	if len(h) < 7 || h[:7] != "Bearer " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
		return
	}
	token := h[7:] // 去掉 "Bearer " 前缀

	if err := s.Auth.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
