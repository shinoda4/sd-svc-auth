package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	stdhttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shinoda4/sd-svc-auth/internal/dto"
	"github.com/shinoda4/sd-svc-auth/internal/repo"
)

func (s *Server) HandleRegister(c *gin.Context) {

	sendEmail := c.DefaultQuery("sendEmail", "true")
	sendEmailBool := sendEmail == "true"

	var body dto.RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(stdhttp.StatusBadRequest, gin.H{
				"error": "empty request body",
			})
			return
		}

		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(stdhttp.StatusConflict, gin.H{"error": "email already exists", "details": e.Email})
		} else {
			c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": "register failed", "details": err.Error()})
		}
		return
	}

	c.JSON(stdhttp.StatusCreated, dto.RegisterResponse{
		Message:     "registered",
		VerifyToken: verifyToken,
	})
}

func (s *Server) HandleLogin(c *gin.Context) {
	var body dto.LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": "empty request body"})
			return
		}
		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accessToken, refreshToken, accessTTL, refreshTTL, err := s.Auth.Login(ctx, body.Email, body.Password)
	if err != nil {
		c.JSON(stdhttp.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(stdhttp.StatusOK, dto.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        int(accessTTL.Seconds()),
		RefreshExpiresIn: int(refreshTTL.Seconds()),
	})
}

func (s *Server) HandleLogout(c *gin.Context) {
	h := c.GetHeader("Authorization")
	if len(h) < 7 || h[:7] != "Bearer " {
		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": "invalid token format"})
		return
	}
	token := h[7:] // 去掉 "Bearer " 前缀

	if err := s.Auth.Logout(c.Request.Context(), token); err != nil {
		c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(stdhttp.StatusOK, gin.H{"message": "logout successful"})
}
