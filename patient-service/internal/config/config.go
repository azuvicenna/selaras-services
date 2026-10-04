package config

import (
	"net"
	"net/url"
	"os"
)

type Config struct {
	DatabaseURL string
	GRPCPort    string

	StorageEndpoint  string
	StorageAccessKey string
	StorageSecretKey string
	StorageBucket    string
	StorageUseSSL    bool
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

		StorageEndpoint:  getEnv("STORAGE_ENDPOINT", "localhost:9000"),
		StorageAccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
		StorageSecretKey: os.Getenv("STORAGE_SECRET_KEY"),
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