package middlewares

import (
	"sync"
	"time"

	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/utils"
	"github.com/gin-gonic/gin"
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

// initRateLimiter initializes the global rate limiter
func initRateLimiter() {
	once.Do(func() {
		// Default: 100 requests per second with burst of 200
		// Adjust these values based on your needs
		globalLimiter = NewIPRateLimiter(100, 200)

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

// RateLimitMiddleware creates a rate limiting middleware
//
// Limits requests per IP address to prevent abuse and DDoS attacks.
// Uses token bucket algorithm for smooth rate limiting.
//
// Default limits: 100 requests/second with burst of 200
func RateLimitMiddleware() gin.HandlerFunc {
	initRateLimiter()

	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		// Get limiter for this IP
		limiter := globalLimiter.GetLimiter(ip)

		// Check if request is allowed
		if !limiter.Allow() {
			logger.Warnf("rate limit exceeded for IP: %s", ip)
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
