// Package config loads the API settings from environment variables.
package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	Port           string
	DatabaseURL    string
	AllowedOrigins []string
	// Timezone is the campus timezone used to decide whether a session is in the past.
	Timezone string
}

func Load() (Config, error) {
	cfg := Config{
		Port:           getEnv("PORT", "3000"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AllowedOrigins: splitList(getEnv("ALLOWED_ORIGINS", "http://localhost:5500,http://127.0.0.1:5500")),
		Timezone:       getEnv("TIMEZONE", "America/Denver"),
	}
	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
