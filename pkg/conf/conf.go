package conf

import (
	"profile_service/pkg/env"
	"strconv"
)

// Конфиги приложения из енв
type Config struct {
	Port                  string
	AuthServiceHost       string
	AuthServicePort       string
	AuthServiceRetries    int
	AuthServiceRetryDelay int

	SecretKeyAccess  string
	SecretKeyRefresh string
}

func New() *Config {
	return &Config{
		Port:                  env.GetEnv("PORT", "8000"),
		AuthServiceHost:       env.GetEnv("AUTH_SERVICE_HOST", "http://localhost"),
		AuthServicePort:       env.GetEnv("AUTH_SERVICE_PORT", "8080"),
		AuthServiceRetries:    getenvInt(env.GetEnv("AUTH_SERVICE_RETRIES", "3")),
		AuthServiceRetryDelay: getenvInt(env.GetEnv("AUTH_SERVICE_RETRY_DELAY", "500")),
		SecretKeyAccess:       env.GetEnv("ACCESS_SECRET", "secretsecret"),
		SecretKeyRefresh:      env.GetEnv("REFRESH_SECRET", "epstein didn't kill himself"),
	}
}

func (c *Config) GetAuthServiceAddr() string {
	return c.AuthServiceHost + ":" + c.AuthServicePort
}

func getenvInt(key string) int {
	v, err := strconv.Atoi(key)
	if err != nil {
		return 0
	}
	return v
}
