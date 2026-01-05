package config

import (
	"os"
)

type Config struct {
	Port            string
	DatabaseURL     string
	GithubClientID  string
	GithubSecret    string
	JWTSecret       string
	FrontendURL     string
	Environment     string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		GithubClientID:  getEnv("GITHUB_CLIENT_ID", ""),
		GithubSecret:    getEnv("GITHUB_CLIENT_SECRET", ""),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		FrontendURL:     getEnv("FRONTEND_URL", "http://localhost:3000"),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
