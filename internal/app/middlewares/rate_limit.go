package middlewares

import (
	"fmt"
	"sync"
	"time"

	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// IPRateLimiter manages rate limiters per IP address
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit // requests per second
	b   int        // burst size
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// GetLimiter returns the rate limiter for the given IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// CleanupOldLimiters removes limiters that haven't been used recently
// Call this periodically to prevent memory leaks
func (i *IPRateLimiter) CleanupOldLimiters() {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Remove limiters with full burst capacity (haven't been used)
	for ip, limiter := range i.ips {
		if limiter.Tokens() == float64(i.b) {
			delete(i.ips, ip)
		}
	}
}

var (
	// Global rate limiter instance
	globalLimiter *IPRateLimiter
	once          sync.Once
)

// initRateLimiter initializes the global rate limiter from RATE_LIMIT_RPS and RATE_LIMIT_BURST env.
func initRateLimiter() {
	once.Do(func() {
		rps := viper.GetInt("RATE_LIMIT_RPS")
		if rps <= 0 {
			rps = 100
		}
		burst := viper.GetInt("RATE_LIMIT_BURST")
		if burst <= 0 {
			burst = 200
		}
		globalLimiter = NewIPRateLimiter(rate.Limit(rps), burst)

		// Cleanup old limiters every 5 minutes
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				globalLimiter.CleanupOldLimiters()
				logger.Debugf("rate limiter cleanup completed")
			}
		}()
	})
}

// rateLimitKey returns the key for rate limiting: user ID if RATE_LIMIT_USE_USER is true
// and user_id is set in context (e.g. by AuthMiddleware), otherwise client IP.
func rateLimitKey(c *gin.Context) string {
	if viper.GetBool("RATE_LIMIT_USE_USER") {
		if userID, exists := c.Get("user_id"); exists {
			if id, ok := userID.(uint); ok {
				return fmt.Sprintf("user:%d", id)
			}
		}
	}
	return c.ClientIP()
}

// RateLimitMiddleware creates a rate limiting middleware.
//
// Limits requests per IP (or per user if RATE_LIMIT_USE_USER=true and user_id is in context).
// Reads RATE_LIMIT_RPS and RATE_LIMIT_BURST from config; if unset or <= 0, uses defaults 100 req/s and burst 200.
func RateLimitMiddleware() gin.HandlerFunc {
	initRateLimiter()

	return func(c *gin.Context) {
		key := rateLimitKey(c)
		limiter := globalLimiter.GetLimiter(key)

		if !limiter.Allow() {
			logger.Warnf("rate limit exceeded for key: %s", key)
			utils.TooManyRequests(c, nil, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddlewareWithConfig creates a rate limiting middleware with custom config
//
// Parameters:
//   - requestsPerSecond: Maximum requests allowed per second
//   - burst: Maximum burst size (allows temporary spikes)
//
// Example:
//
//	router.Use(middlewares.RateLimitMiddlewareWithConfig(50, 100))
func RateLimitMiddlewareWithConfig(requestsPerSecond int, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(requestsPerSecond), burst)

	// Cleanup periodically
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.CleanupOldLimiters()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		ipLimiter := limiter.GetLimiter(ip)

		if !ipLimiter.Allow() {
			logger.Warnf("rate limit exceeded for IP: %s (limit: %d/s, burst: %d)", ip, requestsPerSecond, burst)
			utils.TooManyRequests(c, nil, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}
