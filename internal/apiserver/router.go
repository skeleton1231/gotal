package apiserver

import (
	"github.com/gin-gonic/gin"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {

	g.GET("/api-test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is Test API",
		})
	})
	return g
}
