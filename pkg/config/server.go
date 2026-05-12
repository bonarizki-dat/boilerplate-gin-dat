package config

import (
	"fmt"

	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
)

type ServerConfiguration struct {
	Host                   string
	Port                   string
	JWTSecret              string
	Debug                  bool
	AllowedHosts           string
	LimitCountPerRequest   int64
	Timezone               string // IANA name, e.g. UTC, Asia/Jakarta
	RequestTimeoutSeconds  int    // http.Server read/write timeout
	ShutdownTimeoutSeconds int    // graceful shutdown max wait
	// Rate limit (used by middlewares)
	RateLimitRPS     int
	RateLimitBurst   int
	RateLimitUseUser bool
}

// ServerConfig returns the server address (host:port) from loaded config.
// Call after SetupConfig(). Returns default "0.0.0.0:8000" if config not loaded.
func ServerConfig() string {
	if c := Get(); c != nil {
		addr := fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
		logger.Infof("server listening at %s", addr)
		return addr
	}
	logger.Infof("server listening at 0.0.0.0:8000 (config not loaded)")
	return "0.0.0.0:8000"
}
