package config

import (
	"net"
	"net/url"
	"os"
)

type Config struct {
	DatabaseURL string
	GRPCPort    string
}

func Load() Config {
	dbURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
		Host:     net.JoinHostPort(getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432")),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: "sslmode=disable",
	}
	return Config{
		DatabaseURL: dbURL.String(),
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}