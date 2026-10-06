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

	StorageEndpoint  string
	StorageAccessKey string
	StorageSecretKey string
	StorageBucket    string
	StorageUseSSL    bool
}

func Load() Config {
	dbURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(getEnv("DB_USER", "postgres"), getEnv("DB_PASSWORD", "postgres")),
		Host:     net.JoinHostPort(getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432")),
		Path:     getEnv("DB_NAME", "patient_db"),
		RawQuery: "sslmode=disable",
	}
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", dbURL.String()),
		DBMaxConns:  getEnvInt32("DB_MAX_CONNS", 25),
		DBMinConns:  getEnvInt32("DB_MIN_CONNS", 5),
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		MetricsPort: getEnv("METRICS_PORT", "9092"),

		StorageEndpoint:  getEnv("STORAGE_ENDPOINT", "localhost:9000"),
		StorageAccessKey: getEnv("STORAGE_ACCESS_KEY", "minioadmin"),
		StorageSecretKey: getEnv("STORAGE_SECRET_KEY", "minioadmin"),
		StorageBucket:    getEnv("STORAGE_BUCKET", "patient-documents"),
		StorageUseSSL:    os.Getenv("STORAGE_USE_SSL") == "true",
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
