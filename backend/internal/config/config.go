package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerPort        string
	DatabaseURL       string
	ClickHouseURL     string
	ClickHouseDB      string
	ClickHouseUser    string
	ClickHousePass    string
	MistralAPIKey     string
	MistralModel      string
	MistralEndpoint   string
	MistralTimeout    time.Duration
	PublicAvatarHosts []string
}

func Load() Config {
	return Config{
		ServerPort:        env("SERVER_PORT", "8080"),
		DatabaseURL:       databaseURL(),
		ClickHouseURL:     env("CLICKHOUSE_URL", "http://localhost:8123"),
		ClickHouseDB:      env("CLICKHOUSE_DATABASE", "recap"),
		ClickHouseUser:    env("CLICKHOUSE_USER", "default"),
		ClickHousePass:    os.Getenv("CLICKHOUSE_PASSWORD"),
		MistralAPIKey:     os.Getenv("MISTRAL_API_KEY"),
		MistralModel:      env("MISTRAL_MODEL", "mistral-small-latest"),
		MistralEndpoint:   os.Getenv("MISTRAL_ENDPOINT"),
		MistralTimeout:    durationEnv("MISTRAL_TIMEOUT", 3*time.Second),
		PublicAvatarHosts: csvEnv("PUBLIC_AVATAR_HOSTS"),
	}
}

func databaseURL() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		env("DB_USER", "avito"),
		env("DB_PASSWORD", "avito_password"),
		env("DB_HOST", "localhost"),
		env("DB_PORT", "5432"),
		env("DB_NAME", "avito_recap"),
		env("DB_SSLMODE", "disable"),
	)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if parsed, err := time.ParseDuration(value); err == nil {
		return parsed
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return fallback
}

func csvEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.ToLower(strings.TrimSpace(part))
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
