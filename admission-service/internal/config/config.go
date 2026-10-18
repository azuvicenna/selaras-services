package config

import (
	"net"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	DBMaxConns  int32
	DBMinConns  int32
	GRPCPort    string
	MetricsPort string
	RabbitMQURL string
	MasterServiceAddr string
}

func Load() Config {
	dbURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(getEnv("DB_USER", "admission"), getEnv("DB_PASSWORD", "change-me")),
		Host:     net.JoinHostPort(getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432")),
		Path:     getEnv("DB_NAME", "admission_db"),
		RawQuery: "sslmode=disable",
	}

	return Config{
		DatabaseURL:       getEnv("DATABASE_URL", dbURL.String()),
		DBMaxConns:        getEnvInt32("DB_MAX_CONNS", 25),
		DBMinConns:        getEnvInt32("DB_MIN_CONNS", 5),
		GRPCPort:          getEnv("GRPC_PORT", "50055"),
		MetricsPort:       getEnv("METRICS_PORT", "9096"),
		RabbitMQURL:       getEnv("RABBITMQ_URL", ""),
		MasterServiceAddr: os.Getenv("MASTER_SERVICE_ADDR"),
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
