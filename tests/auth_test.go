package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shinoda4/sd-svc-auth/internal/handler"
	"github.com/shinoda4/sd-svc-auth/internal/repo"
	"github.com/shinoda4/sd-svc-auth/internal/service"
)

func setupTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db := repo.NewMockUserRepo() // 下面会定义 mock
	cache := repo.NewMockRedis()
	authService := service.NewAuthService(db, cache)
	s := handler.NewServer(authService)

	r := gin.Default()
	r.POST("/register", s.HandleRegister)
	r.POST("/login", s.HandleLogin)
	return r
}

func TestRegister(t *testing.T) {
	server := setupTestServer()

	registerPayload := map[string]string{
		"email":    "test@example.com",
		"password": "123456",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("register failed: code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestLogin(t *testing.T) {
	server := setupTestServer()
	registerPayload := map[string]string{
		"email":    "test@example.com",
		"password": "123456",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("register failed: code=%d body=%s", resp.Code, resp.Body.String())
	}
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "123456",
	}
	body2, _ := json.Marshal(loginPayload)
	req2, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")

	resp2 := httptest.NewRecorder()
	server.ServeHTTP(resp2, req2)

	if resp2.Code != http.StatusOK {
		t.Fatalf("login failed: code=%d body=%s", resp2.Code, resp2.Body.String())
	}

	var result map[string]string
	_ = json.Unmarshal(resp2.Body.Bytes(), &result)
	if result["token"] == "" {
		t.Fatalf("expected token, got empty")
	}
}
