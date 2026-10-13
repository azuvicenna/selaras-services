package config

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL      string
	DBMaxConns       int32
	DBMinConns       int32
	GRPCPort         string
	MetricsPort      string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	CaptchaExpiry    time.Duration
}

func Load() Config {
	dbURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(getEnv("DB_USER", "auth"), getEnv("DB_PASSWORD", "change-me")),
		Host:     net.JoinHostPort(getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432")),
		Path:     getEnv("DB_NAME", "auth_db"),
		RawQuery: "sslmode=disable",
	}

	return Config{
		DatabaseURL:      getEnv("DATABASE_URL", dbURL.String()),
		DBMaxConns:       getEnvInt32("DB_MAX_CONNS", 25),
		DBMinConns:       getEnvInt32("DB_MIN_CONNS", 5),
		GRPCPort:         getEnv("GRPC_PORT", "50057"),
		MetricsPort:      getEnv("METRICS_PORT", "9098"),
		JWTSecret:        getEnv("JWT_SECRET", "selaras-hospital-jwt-secret-key-32bytes!"),
		JWTAccessExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		CaptchaExpiry:    getEnvDuration("CAPTCHA_EXPIRY", 5*time.Minute),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt32(key string, fallback int32) int32 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			return int32(n)
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
