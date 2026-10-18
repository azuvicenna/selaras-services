package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort    string
	MetricsPort string

	// RabbitMQ
	RabbitMQURL        string
	RabbitMQQueue      string
	RabbitMQExchange   string
	RabbitMQRoutingKey string
	RabbitMQDLQ        string
	RabbitMQPrefetch   int
	RabbitMQMaxRetries int

	// SATUSEHAT Platform
	SatusehatAuthURL      string
	SatusehatBaseURL      string
	SatusehatClientID     string
	SatusehatClientSecret string
	SatusehatOrgID        string
	RequestTimeout        time.Duration

	// EMR gRPC Target Service
	EmrServiceAddr    string
	MasterServiceAddr string

	// Concurrency
	WorkerConcurrency int
}

func Load() *Config {
	timeoutSec := getEnvAsInt("REQUEST_TIMEOUT_SEC", 30)

	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8085"),
		MetricsPort:        getEnv("METRICS_PORT", "9095"),
		RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://selaras:selaras123@localhost:5672/"),
		RabbitMQQueue:      getEnv("RABBITMQ_QUEUE", "satusehat.sync"),
		RabbitMQExchange:   getEnv("RABBITMQ_EXCHANGE", "satusehat.exchange"),
		RabbitMQRoutingKey: getEnv("RABBITMQ_ROUTING_KEY", "satusehat.sync.#"),
		RabbitMQDLQ:        getEnv("RABBITMQ_DLQ", "satusehat.sync.dlq"),
		RabbitMQPrefetch:   getEnvAsInt("RABBITMQ_PREFETCH", 10),
		RabbitMQMaxRetries: getEnvAsInt("RABBITMQ_MAX_RETRIES", 3),

		SatusehatAuthURL:      getEnv("SATUSEHAT_AUTH_URL", "https://api-satusehat-stg.dto.kemkes.go.id/oauth2/v1"),
		SatusehatBaseURL:      getEnv("SATUSEHAT_BASE_URL", "https://api-satusehat-stg.dto.kemkes.go.id/fhir-r4/v1"),
		SatusehatClientID:     getEnv("SATUSEHAT_CLIENT_ID", "satusehat-demo-client-id"),
		SatusehatClientSecret: getEnv("SATUSEHAT_CLIENT_SECRET", "satusehat-demo-client-secret"),
		SatusehatOrgID:        getEnv("SATUSEHAT_ORG_ID", "10000004"),
		RequestTimeout:        time.Duration(timeoutSec) * time.Second,

		EmrServiceAddr:    getEnv("EMR_SERVICE_ADDR", "localhost:50052"),
		MasterServiceAddr: getEnv("MASTER_SERVICE_ADDR", ""),
		WorkerConcurrency: getEnvAsInt("WORKER_CONCURRENCY", 5),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
