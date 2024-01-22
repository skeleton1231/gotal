package apiserver

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/internal/apiserver/controller/v1/user"
	"github.com/skeleton1231/gotal/internal/apiserver/store/database"
	"github.com/skeleton1231/gotal/internal/pkg/code"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/internal/pkg/middleware"
	"github.com/skeleton1231/gotal/internal/pkg/middleware/auth"
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
	jwtStrategy, _ := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	// Refresh time can be longer than token timeout
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		response.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	storeIns, _ := database.GetMySQLFactoryOr(nil)
	userController := user.NewUserController(storeIns)
	testController(g)

	authGroup := g.Group("/v1")
	authGroup.Use(auto.AuthFunc())
	{
		// user RESTful resource
		userv1 := authGroup.Group("/users")
		{
			userv1.PUT("/:id", userController.Update)
			userv1.DELETE("/:id", userController.Delete)
		}
	}

	noAuthGroup := g.Group("/v1")
	{
		noAuthGroup.POST("/users", userController.Create)
		noAuthGroup.GET("/users/:id", userController.Get)

	}

	return g
}

func testController(g *gin.Engine) {

	if gin.Mode() != "debug" {
		return
	}

	// // Apply the rate limiter middleware with parameters from Viper
	// g.Use(middleware.RateLimiter())

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
		randomNumber := strconv.Itoa(rand.Int()) // 生成一个随机整数

		redisClient := &cache.RedisClusterV2{}
		// Set the value in Redis
		err := redisClient.SetRawKey(ctx, key, randomNumber, 3600*time.Second)
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
