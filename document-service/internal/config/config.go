package config

import (
	"net"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL   string
	DBMaxConns    int32
	DBMinConns    int32
	GRPCPort      string
	MetricsPort   string
	TypstBinPath  string
	VerifyBaseURL string
}

func Load() Config {
	dbURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(getEnv("DB_USER", "document"), getEnv("DB_PASSWORD", "change-me")),
		Host:     net.JoinHostPort(getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432")),
		Path:     getEnv("DB_NAME", "document_db"),
		RawQuery: "sslmode=disable",
	}

	return Config{
		DatabaseURL:   getEnv("DATABASE_URL", dbURL.String()),
		DBMaxConns:    getEnvInt32("DB_MAX_CONNS", 25),
		DBMinConns:    getEnvInt32("DB_MIN_CONNS", 5),
		GRPCPort:      getEnv("GRPC_PORT", "50062"),
		MetricsPort:   getEnv("METRICS_PORT", "9103"),
		TypstBinPath:  getEnv("TYPST_BIN_PATH", "typst"),
		VerifyBaseURL: getEnv("VERIFY_BASE_URL", "https://selaras.hospital.id/verify"),
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
