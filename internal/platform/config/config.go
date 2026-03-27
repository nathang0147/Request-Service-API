package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	defaultPort             = "8080"
	defaultLogLevel         = "info"
	defaultProvider         = "walt"
	defaultWaltVerifierMode = "legacy"
	databaseURLEnv          = "DATABASE_URL"
	waltVerifierBaseURLEnv  = "WALT_VERIFIER_BASE_URL"
	waltVerifierModeEnv     = "WALT_VERIFIER_MODE"
	waltBearerTokenEnv      = "WALT_BEARER_TOKEN"
	waltVCPolicyWebhookEnv  = "WALT_VC_POLICY_WEBHOOK_URL"
	callbackBaseURLEnv      = "CALLBACK_BASE_URL"
	publicBaseURLEnv        = "PUBLIC_BASE_URL"
	publicRedirectTemplate  = "PUBLIC_REDIRECT_URL_TEMPLATE"
	callbackAuthSecretEnv   = "CALLBACK_AUTH_SECRET"
	portEnv                 = "PORT"
	logLevelEnv             = "LOG_LEVEL"
	defaultProviderEnv      = "DEFAULT_PROVIDER"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	WaltVerifierBaseURL    string
	WaltVerifierMode       string
	WaltBearerToken        string
	WaltVCPolicyWebhookURL string
	CallbackBaseURL        string
	PublicBaseURL          string
	PublicRedirectTemplate string
	LogLevel               string
	DefaultProvider        string
	CallbackAuthSecret     string
}

func Load() (Config, error) {
	var missing []string

	cfg := Config{
		Port:                   envOrDefault(portEnv, defaultPort),
		DatabaseURL:            requiredEnv(databaseURLEnv, &missing),
		WaltVerifierBaseURL:    requiredEnv(waltVerifierBaseURLEnv, &missing),
		WaltVerifierMode:       envOrDefault(waltVerifierModeEnv, defaultWaltVerifierMode),
		WaltBearerToken:        envOrDefault(waltBearerTokenEnv, ""),
		WaltVCPolicyWebhookURL: envOrDefault(waltVCPolicyWebhookEnv, ""),
		CallbackBaseURL:        requiredEnv(callbackBaseURLEnv, &missing),
		PublicBaseURL:          envOrDefault(publicBaseURLEnv, ""),
		PublicRedirectTemplate: envOrDefault(publicRedirectTemplate, ""),
		LogLevel:               envOrDefault(logLevelEnv, defaultLogLevel),
		DefaultProvider:        envOrDefault(defaultProviderEnv, defaultProvider),
		CallbackAuthSecret:     requiredEnv(callbackAuthSecretEnv, &missing),
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
