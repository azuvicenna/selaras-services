package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort    string
	MetricsPort string

	// RabbitMQ Message Broker
	RabbitMQURL        string
	RabbitMQQueue      string
	RabbitMQExchange   string
	RabbitMQRoutingKey string
	RabbitMQDLQ        string
	RabbitMQPrefetch   int
	RabbitMQMaxRetries int

	// BPJS Web Services (TrustMark / VClaim / Antrean / Aplicares)
	BPJSConsID          string
	BPJSSecretKey       string
	BPJSUserKey         string
	BPJSVClaimBaseURL   string
	BPJSAntreanBaseURL  string
	BPJSAplicaresBaseURL string
	BPJSKodePPK         string
	RequestTimeout      time.Duration

	// Downstream gRPC Services
	EmrServiceAddr    string
	BedServiceAddr    string
	MasterServiceAddr string

	// Concurrency
	WorkerConcurrency int
}

func Load() *Config {
	timeoutSec := getEnvAsInt("REQUEST_TIMEOUT_SEC", 30)

	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8086"),
		MetricsPort:        getEnv("METRICS_PORT", "9096"),
		RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://selaras:selaras123@localhost:5672/"),
		RabbitMQQueue:      getEnv("RABBITMQ_QUEUE", "bpjs.sync"),
		RabbitMQExchange:   getEnv("RABBITMQ_EXCHANGE", "bpjs.exchange"),
		RabbitMQRoutingKey: getEnv("RABBITMQ_ROUTING_KEY", "bpjs.#"),
		RabbitMQDLQ:        getEnv("RABBITMQ_DLQ", "bpjs.sync.dlq"),
		RabbitMQPrefetch:   getEnvAsInt("RABBITMQ_PREFETCH", 10),
		RabbitMQMaxRetries: getEnvAsInt("RABBITMQ_MAX_RETRIES", 3),

		BPJSConsID:           getEnv("BPJS_CONS_ID", "12345"),
		BPJSSecretKey:        getEnv("BPJS_SECRET_KEY", "bpjs-secret-key-32chars-example"),
		BPJSUserKey:          getEnv("BPJS_USER_KEY", "bpjs-user-key-trustmark-dev"),
		BPJSVClaimBaseURL:    getEnv("BPJS_VCLAIM_BASE_URL", "https://apijkn-dev.bpjs-kesehatan.go.id/vclaim-rest-dev"),
		BPJSAntreanBaseURL:   getEnv("BPJS_ANTREAN_BASE_URL", "https://apijkn-dev.bpjs-kesehatan.go.id/antreanrs_dev"),
		BPJSAplicaresBaseURL: getEnv("BPJS_APLICARES_BASE_URL", "https://apijkn-dev.bpjs-kesehatan.go.id/aplicaresws"),
		BPJSKodePPK:          getEnv("BPJS_KODE_PPK", "0115R001"),
		RequestTimeout:       time.Duration(timeoutSec) * time.Second,

		EmrServiceAddr:    getEnv("EMR_SERVICE_ADDR", "localhost:50052"),
		BedServiceAddr:    getEnv("BED_SERVICE_ADDR", "localhost:50054"),
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
