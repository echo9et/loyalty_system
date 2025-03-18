package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gophermart.ru/internal/entities"
	"gophermart.ru/internal/utils"
	"gophermart.ru/mocks"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserManagment := mocks.NewMockUserManagment(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routers := router.Group("/api/user")
	Login(routers.Group("/login"), mockUserManagment)

	type SendUser struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	t.Run("Test successful login", func(t *testing.T) {
		user := SendUser{
			Login:    "test",
			Password: "password",
		}
		hashedPassword := utils.Sha256hash(user.Password)

		storedUser := entities.User{
			Login:        user.Login,
			HashPassword: hashedPassword,
		}

		mockUserManagment.EXPECT().
			User(user.Login).
			Return(&storedUser, nil)

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"result":"success"`)
	})

	t.Run("Test empty fields", func(t *testing.T) {
		user := SendUser{
			Login:    "",
			Password: "",
		}

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"неверный формат запроса"`)
	})

	t.Run("Test invalid credentials", func(t *testing.T) {
		user := SendUser{
			Login:    "test",
			Password: "wrongpassword",
		}
		hashedPassword := utils.Sha256hash("correctpassword")

		storedUser := entities.User{
			Login:        user.Login,
			HashPassword: hashedPassword,
		}

		mockUserManagment.EXPECT().
			User(user.Login).
			Return(&storedUser, nil)

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"неверная пара логин/пароль"`)
	})

	t.Run("Test user not found", func(t *testing.T) {
		user := SendUser{
			Login:    "nonexistent",
			Password: "password",
		}

		mockUserManagment.EXPECT().
			User(user.Login).
			Return(nil, fmt.Errorf("user not found"))

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"внутренняя ошибка сервера"`)
	})
}
