package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	sync.Mutex
	attempts map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.Lock()
		defer rl.Unlock()

		now := time.Now()
		windowStart := now.Add(-rl.window)

		// Remove old attempts
		var validAttempts []time.Time
		for _, attempt := range rl.attempts[ip] {
			if attempt.After(windowStart) {
				validAttempts = append(validAttempts, attempt)
			}
		}
		rl.attempts[ip] = validAttempts

		// Check if limit exceeded
		if len(validAttempts) >= rl.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		// Add current attempt
		rl.attempts[ip] = append(rl.attempts[ip], now)
		c.Next()
	}
}
