package interfaces

import "github.com/gin-gonic/gin"

type Server struct {
}

func New() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		if c.Request.Method != "POST" || c.Request.Method != "GET" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
