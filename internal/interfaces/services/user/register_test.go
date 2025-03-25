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

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserManagment := mocks.NewMockUserManagment(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routers := router.Group("/api/user")
	Register(routers.Group("/register"), mockUserManagment)

	type SendUser struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	t.Run("Test successful registration", func(t *testing.T) {
		user := SendUser{
			Login:    "test",
			Password: "password",
		}
		hashedPassword := utils.Sha256hash(user.Password)

		mockUserManagment.EXPECT().
			InsertUser(entities.User{
				Login:        user.Login,
				HashPassword: hashedPassword,
			}).
			Return(nil)

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test empty fields", func(t *testing.T) {
		user := SendUser{
			Login:    "",
			Password: "",
		}

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test login already exists", func(t *testing.T) {
		user := SendUser{
			Login:    "test",
			Password: "password",
		}
		hashedPassword := utils.Sha256hash(user.Password)

		mockUserManagment.EXPECT().
			InsertUser(entities.User{
				Login:        user.Login,
				HashPassword: hashedPassword,
			}).
			Return(fmt.Errorf("логин уже занят"))

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
