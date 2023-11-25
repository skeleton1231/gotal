package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/pkg/log"
)

func ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在请求被处理之前
		c.Next()

		// 在请求被处理之后
		// 获取响应状态码和数据
		statusCode := c.Writer.Status()
		data, exists := c.Get("response")
		if exists {
			log.Debugf("Request processed: status=%d, data=%v", statusCode, data)
		}
	}
}
