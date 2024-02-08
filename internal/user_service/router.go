package userservice

import (
	"net/http"
	"time"

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

	// Middlewares.

	return g
}

func testController(g *gin.Engine) {

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
		key := "testKey"
		// randomNumber := strconv.Itoa(rand.Int()) // 生成一个随机整数

		redisClient := &cache.RedisClusterV2{}
		// Set the value in Redis
		err := redisClient.SetRawKey(ctx, key, "123", 3600*time.Second)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to set key in Redis",
			})
			return
		}

		// Get the value back from Redis
		retrievedValue, err := redisClient.GetRawKey(ctx, key)
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
