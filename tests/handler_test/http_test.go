package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shinoda4/sd-svc-auth/tests/testserver"
)

func TestRegister(t *testing.T) {
	server := testserver.SetupFullTestServer()

	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"username": "test",
		"password": "123456",
	})
	req, _ := http.NewRequest("POST", "/api/v1/register?sendEmail=false", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("register failed: code=%d body=%s", resp.Code, resp.Body.String())
	}
}
func TestLogin(t *testing.T) {
	server := testserver.SetupFullTestServer()

	// 先注册
	bodyReg, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"username": "test",
		"password": "123456",
	})
	reqReg, _ := http.NewRequest("POST", "/api/v1/register?sendEmail=false", bytes.NewBuffer(bodyReg))
	reqReg.Header.Set("Content-Type", "application/json")
	respReg := httptest.NewRecorder()
	server.ServeHTTP(respReg, reqReg)

	if respReg.Code != http.StatusCreated {
		t.Fatalf("register failed: %s", respReg.Body.String())
	}

	// 登录
	bodyLogin, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "123456",
	})
	reqLogin, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(bodyLogin))
	reqLogin.Header.Set("Content-Type", "application/json")

	respLogin := httptest.NewRecorder()
	server.ServeHTTP(respLogin, reqLogin)

	if respLogin.Code != http.StatusOK {
		t.Fatalf("login failed: %s", respLogin.Body.String())
	}
}

func TestRefresh(t *testing.T) {
	server := testserver.SetupFullTestServer()

	// 注册并登录
	registerPayload := map[string]string{"email": "test@example.com", "username": "test", "password": "123456"}
	body, _ := json.Marshal(registerPayload)
	server.ServeHTTP(httptest.NewRecorder(),
		mustReq("POST", "/api/v1/register?sendEmail=false", body))

	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, mustReq("POST", "/api/v1/login", body))

	var reply map[string]any
	_ = json.Unmarshal(resp.Body.Bytes(), &reply)
	refreshToken := reply["refresh_token"].(string)

	// refresh
	refPayload := map[string]string{"refresh_token": refreshToken}
	body2, _ := json.Marshal(refPayload)

	resp2 := httptest.NewRecorder()
	server.ServeHTTP(resp2, mustReq("POST", "/api/v1/refresh", body2))

	if resp2.Code != http.StatusOK {
		t.Fatalf("refresh failed: %s", resp2.Body.String())
	}
}

func TestVerifyToken(t *testing.T) {
	server := testserver.SetupFullTestServer()

	// 注册 + 登录
	reg := map[string]string{"email": "test@example.com", "username": "test", "password": "123456"}
	body, _ := json.Marshal(reg)
	server.ServeHTTP(httptest.NewRecorder(), mustReq("POST", "/api/v1/register?sendEmail=false", body))

	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, mustReq("POST", "/api/v1/login", body))

	var login map[string]any
	_ = json.Unmarshal(resp.Body.Bytes(), &login)
	access := login["access_token"].(string)

	verifyPayload := map[string]string{"token": access}
	vbody, _ := json.Marshal(verifyPayload)

	resp2 := httptest.NewRecorder()
	server.ServeHTTP(resp2, mustReq("POST", "/api/v1/verify-token", vbody))

	if resp2.Code != http.StatusOK {
		t.Fatalf("verify failed: %s", resp2.Body.String())
	}
}

func TestMe(t *testing.T) {
	server := testserver.SetupFullTestServer()

	// 注册 + 登录
	reg := map[string]string{"email": "aaa@example.com", "username": "test", "password": "123456"}
	body, _ := json.Marshal(reg)
	server.ServeHTTP(httptest.NewRecorder(), mustReq("POST", "/api/v1/register?sendEmail=false", body))

	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, mustReq("POST", "/api/v1/login", body))

	var login map[string]any
	_ = json.Unmarshal(resp.Body.Bytes(), &login)
	access := login["access_token"].(string)

	req := mustReq("GET", "/api/v1/authorized/me", nil)
	req.Header.Set("Authorization", "Bearer "+access)

	resp2 := httptest.NewRecorder()
	server.ServeHTTP(resp2, req)

	if resp2.Code != http.StatusOK {
		t.Fatalf("/me failed: %s", resp2.Body.String())
	}
}

func TestLogout(t *testing.T) {
	server := testserver.SetupFullTestServer()

	// 注册 + 登录
	reg := map[string]string{"email": "aaa@example.com", "username": "test", "password": "123456"}
	body, _ := json.Marshal(reg)
	server.ServeHTTP(httptest.NewRecorder(), mustReq("POST", "/api/v1/register?sendEmail=false", body))

	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, mustReq("POST", "/api/v1/login", body))

	var login map[string]any
	_ = json.Unmarshal(resp.Body.Bytes(), &login)
	access := login["access_token"].(string)

	// logout
	req := mustReq("POST", "/api/v1/logout", nil)
	req.Header.Set("Authorization", "Bearer "+access)

	resp2 := httptest.NewRecorder()
	server.ServeHTTP(resp2, req)

	if resp2.Code != http.StatusOK {
		t.Fatalf("logout failed: %s", resp2.Body.String())
	}
}

// small helper
func mustReq(method, path string, body []byte) *http.Request {
	var buf *bytes.Buffer
	if body != nil {
		buf = bytes.NewBuffer(body)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}
