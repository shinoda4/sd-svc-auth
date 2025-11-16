package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shinoda4/sd-svc-auth/internal/service/auth"
)

type Server struct {
	Auth *auth.Service
}

func NewServer(auth *auth.Service) *Server {
	return &Server{Auth: auth}
}
func StartServer(authService *auth.Service) {
	s := &Server{Auth: authService}
	port := os.Getenv("SERVER_PORT")

	r := gin.Default()

	api := r.Group("/api/v1")

	api.POST("/register", s.HandleRegister)
	api.POST("/login", s.HandleLogin)
	api.POST("/refresh", s.HandleRefresh)
	api.POST("/verify-token", s.HandleVerifyToken)
	api.POST("/logout", s.HandleLogout)
	api.GET("/verify", s.HandleVerifyEmail)

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
