package testserver

import (
	"github.com/gin-gonic/gin"
	"github.com/shinoda4/sd-svc-auth/internal/handler"
	"github.com/shinoda4/sd-svc-auth/internal/repo"
	"github.com/shinoda4/sd-svc-auth/internal/service"
)

func SetupFullTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db := repo.NewMockUserRepo()
	cache := repo.NewMockRedis()
	authService := service.NewAuthService(db, cache)
	s := handler.NewServer(authService)

	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.POST("/register", s.HandleRegister)
		api.POST("/login", s.HandleLogin)
		api.POST("/refresh", s.HandleRefresh)
		api.POST("/verify-token", s.HandleVerifyToken)
		api.POST("/logout", s.HandleLogout)
		api.GET("/verify", s.HandleVerifyEmail)
	}

	auth := api.Group("/authorized")
	auth.Use(s.JwtMiddleware())
	{
		auth.GET("/me", s.HandleMe)
	}

	return r
}
