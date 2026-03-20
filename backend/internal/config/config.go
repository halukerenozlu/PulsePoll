package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Auth     AuthConfig
}

type AppConfig struct {
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type AuthConfig struct {
	JWTSecret           string
	AccessTokenTTLMin   int
	RefreshTokenTTLHour int
	RefreshCookieName   string
	RefreshCookieSecure bool
}

func Load() Config {
	_ = godotenv.Load(".env", ".env.local", "../.env", "../.env.local")

	return Config{
		App: AppConfig{
			Port: getEnv("PORT", "8080"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "ephemeral"),
			Password: getEnv("POSTGRES_PASSWORD", "ephemeral"),
			DBName:   getEnv("POSTGRES_DB", "ephemeral"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Auth: AuthConfig{
			JWTSecret:           getEnv("AUTH_JWT_SECRET", "change-me-in-production"),
			AccessTokenTTLMin:   getEnvAsInt("AUTH_ACCESS_TOKEN_TTL_MIN", 15),
			RefreshTokenTTLHour: getEnvAsInt("AUTH_REFRESH_TOKEN_TTL_HOUR", 24),
			RefreshCookieName:   getEnv("AUTH_REFRESH_COOKIE_NAME", "refresh_token"),
			RefreshCookieSecure: getEnvAsBool("AUTH_REFRESH_COOKIE_SECURE", false),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvAsBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
