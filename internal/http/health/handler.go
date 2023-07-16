package health

import "github.com/gin-gonic/gin"

func NewServer() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "healthy")
	})

	return r
}
