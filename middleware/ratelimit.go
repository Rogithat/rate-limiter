package middleware

import (
	"net/http"

	"rate-limiter/limiter"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(limiter *limiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		limited, err := limiter.CheckRateLimit(c.Request.Context(), ip, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		if limited {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
