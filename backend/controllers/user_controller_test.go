package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-backend/config"
	"gin-backend/routes"

	"github.com/stretchr/testify/assert"
)

func init() {
	// 初始化配置
	config.LoadConfig()
}

func TestGetUsers(t *testing.T) {
	// 设置路由
	router := routes.SetupRouter(config.DB)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])
	assert.NotNil(t, response["data"])
}

func TestGetUser(t *testing.T) {
	router := routes.SetupRouter(config.DB)

	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.NotNil(t, response["data"])
}

func TestGetUserNotFound(t *testing.T) {
	router := routes.SetupRouter(config.DB)

	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(404), response["code"])
}

func TestHealthCheck(t *testing.T) {
	router := routes.SetupRouter(config.DB)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}
