package user

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	config "gophermart.ru/internal"
)

func TestLogin_Success(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Создание роутера
	router := gin.New()
	group := router.Group("/login")
	Login(group, nil)

	// Настройка конфигурации для теста
	config.Get().AliveToken = time.Minute
	config.Get().SecretKey = "test-secret-key"

	// Создание тестового запроса
	w := httptest.NewRecorder()
	reqBody := `{"username":"testuser","password":"testpassword"}`
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(strings.NewReader(reqBody))

	// Выполнение запроса
	router.ServeHTTP(w, req)

	// Проверка статус-кода
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка содержимого ответа
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response["result"])

	// Проверка куки
	cookie := w.Result().Cookies()
	assert.Len(t, cookie, 1)
	tokenCookie := cookie[0]
	assert.Equal(t, "token", tokenCookie.Name)
	assert.NotEmpty(t, tokenCookie.Value)

	// Проверка токена JWT
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().SecretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, "testuser", claims.Login)
}

func TestLogin_BadRequest_EmptyFields(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Создание роутера
	router := gin.New()
	group := router.Group("/login")
	Login(group, nil)

	// Создание тестового запроса с пустыми полями
	w := httptest.NewRecorder()
	reqBody := `{"username":"","password":""}`
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(strings.NewReader(reqBody))

	// Выполнение запроса
	router.ServeHTTP(w, req)

	// Проверка статус-кода
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Проверка содержимого ответа
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Bad Request", response["error"])
}

func TestLogin_InvalidJSON(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Создание роутера
	router := gin.New()
	group := router.Group("/login")
	Login(group, nil)

	// Создание тестового запроса с некорректным JSON
	w := httptest.NewRecorder()
	reqBody := `invalid-json`
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(strings.NewReader(reqBody))

	// Выполнение запроса
	router.ServeHTTP(w, req)

	// Проверка статус-кода
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Проверка содержимого ответа
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Bad Request", response["error"])
}
