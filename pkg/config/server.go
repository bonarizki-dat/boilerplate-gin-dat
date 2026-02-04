package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Host                 string
	Port                 string
	Secret               string
	JWTSecret            string
	Debug                bool
	AllowedHosts         string
	LimitCountPerRequest int64
	// Rate limit (used by middlewares)
	RateLimitRPS     int
	RateLimitBurst   int
	RateLimitUseUser bool
}

func ServerConfig() string {
	if c := Get(); c != nil {
		addr := fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
		log.Print("Server Running at :", addr)
		return addr
	}
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")
	addr := fmt.Sprintf("%s:%s", viper.GetString("SERVER_HOST"), viper.GetString("SERVER_PORT"))
	log.Print("Server Running at :", addr)
	return addr
}
