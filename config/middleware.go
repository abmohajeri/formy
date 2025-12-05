package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

var limiterMap sync.Map

func RateLimitMiddleware(interval time.Duration, requestsPerInterval int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip, interval, requestsPerInterval)
		if !limiter.Allow() {
			fmt.Printf("Middleware => Rate limit exceeded for IP %s\n", ip)
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

func getLimiter(ip string, interval time.Duration, requestsPerInterval int) *rate.Limiter {
	actual, exists := limiterMap.Load(ip)
	if exists {
		return actual.(*rate.Limiter)
	}

	// Define the rate limit (tokens per second)
	rateLimit := rate.Limit(float64(requestsPerInterval) / interval.Seconds())
	limiter := rate.NewLimiter(rateLimit, requestsPerInterval)
	limiterMap.Store(ip, limiter)

	return limiter
}
