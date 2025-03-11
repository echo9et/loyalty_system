package interfaces

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	config "gophermart.ru/internal"
	"gophermart.ru/internal/infrastructure/storage"
	"gophermart.ru/internal/interfaces/services"
)

type Server struct {
	Engine *gin.Engine
	db     *storage.Database
}

type DataBaser interface {
}

type ServerRouter interface {
}

func New() (*Server, error) {
	server := &Server{}
	db, err := storage.NewDatabase(config.GetAddrDatabase())

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	server.db = db
	gin.SetMode(gin.ReleaseMode)

	server.Engine = gin.New()
	server.Engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		if c.Request.Method != "POST" && c.Request.Method != "GET" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	server.Engine.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.Next()
	})
	server.Engine.Use(MidlewareAuth)
	server.Engine.Use(Midleware2)

	services.User(server.Engine.Group("api/user"), server.db)

	routes := server.Engine.Routes()
	slog.Info("Зарегистрированные маршруты:")
	for _, route := range routes {
		slog.Info("Route",
			"method", route.Method,
			"path", route.Path,
		)
	}

	return server, nil
}
