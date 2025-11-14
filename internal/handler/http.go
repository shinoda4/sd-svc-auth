package handler

import (
	"context"
	"errors"
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

	r.POST("/register", s.HandleRegister)
	r.POST("/login", s.HandleLogin)

	authorized := r.Group("/")
	authorized.Use(s.jwtMiddleware())
	{
		authorized.GET("/me", s.handleMe)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := s.Auth.Register(ctx, body.Email, body.Password); err != nil {
		var e *repo.ErrUserExists
		if errors.As(err, &e) {
			c.JSON(http.StatusConflict, gin.H{"error": "register user exists", "details": e.Email})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registered"})
}

type loginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) HandleLogin(c *gin.Context) {
	var body loginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	token, ttl, err := s.Auth.Login(ctx, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "expires_in": int(ttl.Seconds())})
}

func (s *Server) handleMe(c *gin.Context) {
	claims, _ := c.Get("claims")
	c.JSON(http.StatusOK, gin.H{"claims": claims})
}

func (s *Server) jwtMiddleware() gin.HandlerFunc {
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
