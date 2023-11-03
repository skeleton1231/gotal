package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/pkg/cache"
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

	// Redis Test API endpoint to set and get a value
	g.GET("/redis-test", func(c *gin.Context) {
		ctx := c.Request.Context()

		// Example key-value to set in Redis
		key := "testKey"
		value := "testValue"

		redisClient := cache.RedisClusterV2{KeyPrefix: "test"}
		// Set the value in Redis
		err := redisClient.SetKey(ctx, key, value, 3600)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to set key in Redis",
			})
			return
		}

		// Get the value back from Redis
		retrievedValue, err := redisClient.GetKey(ctx, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get key from Redis",
			})
			return
		}

		// Return the retrieved value
		c.JSON(http.StatusOK, gin.H{
			"retrieved_value": retrievedValue,
		})
	})
	return g
}
