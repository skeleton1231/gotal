package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/internal/pkg/middleware"
	"github.com/skeleton1231/gotal/internal/pkg/response"
	"github.com/skeleton1231/gotal/pkg/cache"
	"github.com/skeleton1231/gotal/pkg/log"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
	g.Use(middleware.ResponseLogger())
}

func installController(g *gin.Engine) *gin.Engine {
	testController(g)
	return g
}

func testController(g *gin.Engine) {

	if gin.Mode() != "debug" {
		return
	}

	// Apply the rate limiter middleware with parameters from Viper
	g.Use(middleware.RateLimiter())

	g.GET("/api-test", func(c *gin.Context) {
		log.Info("Logger testing")
		c.JSON(200, gin.H{
			"message": "This is Test API",
		})
	})

	// Redis Test API endpoint to set and get a value
	g.GET("/redis-test", func(c *gin.Context) {
		ctx := c.Request.Context()

		// Example key-value to set in Redis
		key := "testKey111"
		value := "testValue111"

		redisClient := &cache.RedisClusterV2{KeyPrefix: "test"}
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

	g.GET("/test-response", testHandler)

}

func testHandler(c *gin.Context) {
	// Check for the existence of the 'error' query parameter
	if _, exists := c.GetQuery("error"); exists {
		// Simulating an error scenario
		err := errors.WithCode(code.ErrValidation, "Simulated error message")
		response.WriteResponse(c, err, nil)
	} else {
		// Simulating a success scenario
		response.WriteResponse(c, nil, gin.H{"message": "Test successful"})
	}
}
