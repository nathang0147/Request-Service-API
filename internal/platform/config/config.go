package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	defaultPort           = "8080"
	defaultLogLevel       = "info"
	defaultProvider       = "walt"
	databaseURLEnv        = "DATABASE_URL"
	waltBaseURLEnv        = "WALT_BASE_URL"
	waltAPIKeyEnv         = "WALT_API_KEY"
	callbackBaseURLEnv    = "CALLBACK_BASE_URL"
	callbackAuthSecretEnv = "CALLBACK_AUTH_SECRET"
	portEnv               = "PORT"
	logLevelEnv           = "LOG_LEVEL"
	defaultProviderEnv    = "DEFAULT_PROVIDER"
)

type Config struct {
	Port               string
	DatabaseURL        string
	WaltBaseURL        string
	WaltAPIKey         string
	CallbackBaseURL    string
	LogLevel           string
	DefaultProvider    string
	CallbackAuthSecret string
}

func Load() (Config, error) {
	var missing []string

	cfg := Config{
		Port:               envOrDefault(portEnv, defaultPort),
		DatabaseURL:        requiredEnv(databaseURLEnv, &missing),
		WaltBaseURL:        requiredEnv(waltBaseURLEnv, &missing),
		WaltAPIKey:         requiredEnv(waltAPIKeyEnv, &missing),
		CallbackBaseURL:    requiredEnv(callbackBaseURLEnv, &missing),
		LogLevel:           envOrDefault(logLevelEnv, defaultLogLevel),
		DefaultProvider:    envOrDefault(defaultProviderEnv, defaultProvider),
		CallbackAuthSecret: requiredEnv(callbackAuthSecretEnv, &missing),
	}

	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func requiredEnv(key string, missing *[]string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		*missing = append(*missing, key)
	}

	return value
}
