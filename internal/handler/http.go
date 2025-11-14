package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shinoda4/sd-svc-auth/internal/repo"
	"github.com/shinoda4/sd-svc-auth/internal/service"
)

type Server struct {
	Auth *service.AuthService
}

func NewServer(auth *service.AuthService) *Server {
	return &Server{Auth: auth}
}
func StartServer(authService *service.AuthService) {
	s := &Server{Auth: authService}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	api := r.Group("/api/v1")

	api.POST("/register", s.HandleRegister)
	api.POST("/login", s.HandleLogin)
	api.POST("/refresh", s.HandleRefresh)
	api.POST("/verify", s.HandleVerify)
	api.POST("/logout", s.HandleLogout)

	authorized := api.Group("/authorized")
	authorized.Use(s.JwtMiddleware())
	{
		authorized.GET("/me", s.HandleMe)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	_ = srv.ListenAndServe()
}

type registerBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) HandleRegister(c *gin.Context) {
	var body registerBody
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
	if err := s.Auth.Register(ctx, body.Email, body.Password); err != nil {
		var e *repo.ErrUserExists
		if errors.As(err, &e) {
			c.JSON(http.StatusConflict, gin.H{"error": "register email exists", "details": e.Email})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

type loginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) HandleLogin(c *gin.Context) {
	var body loginBody
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

	c.JSON(http.StatusOK, gin.H{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"expires_in":         int(accessTTL.Seconds()),
		"refresh_expires_in": int(refreshTTL.Seconds()),
	})
}

type refreshBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *Server) HandleRefresh(c *gin.Context) {
	var body refreshBody
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

func (s *Server) HandleMe(c *gin.Context) {
	claims, _ := c.Get("claims")
	c.JSON(http.StatusOK, gin.H{"claims": claims})
}

func (s *Server) HandleVerify(c *gin.Context) {
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

func (s *Server) JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if len(h) < 7 || h[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := h[7:]
		claims, err := s.Auth.ValidateToken(c.Request.Context(), tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
