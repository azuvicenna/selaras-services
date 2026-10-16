package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort    string
	MetricsPort string

	// RabbitMQ Settings
	RabbitMQURL        string
	RabbitMQQueue      string
	RabbitMQExchange   string
	RabbitMQRoutingKey string
	RabbitMQDLQ        string

	// WhatsApp Gateway Settings
	WhatsAppProvider string // "mock", "http"
	WhatsAppAPIURL   string
	WhatsAppAPIToken string

	// Email / SMTP Settings
	EmailProvider string // "mock", "smtp"
	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPassword  string
	SMTPFrom      string
	SMTPFromName  string

	// Worker Concurrency & Request Timeouts
	WorkerConcurrency int
	RequestTimeout    time.Duration
}

func Load() *Config {
	timeoutSec := getEnvAsInt("REQUEST_TIMEOUT_SEC", 30)

	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8087"),
		MetricsPort:        getEnv("METRICS_PORT", "9097"),
		RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://selaras:selaras123@localhost:5672/"),
		RabbitMQQueue:      getEnv("RABBITMQ_QUEUE", "notification.send"),
		RabbitMQExchange:   getEnv("RABBITMQ_EXCHANGE", "notification.exchange"),
		RabbitMQRoutingKey: getEnv("RABBITMQ_ROUTING_KEY", "notification.#"),
		RabbitMQDLQ:        getEnv("RABBITMQ_DLQ", "notification.send.dlq"),

		WhatsAppProvider: getEnv("WHATSAPP_PROVIDER", "mock"),
		WhatsAppAPIURL:   getEnv("WHATSAPP_API_URL", "https://api.fonnte.com/send"),
		WhatsAppAPIToken: getEnv("WHATSAPP_API_TOKEN", "demo-token"),

		EmailProvider: getEnv("EMAIL_PROVIDER", "mock"),
		SMTPHost:      getEnv("SMTP_HOST", "localhost"),
		SMTPPort:      getEnvAsInt("SMTP_PORT", 1025),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:      getEnv("SMTP_FROM", "no-reply@selaras.hospital.id"),
		SMTPFromName:  getEnv("SMTP_FROM_NAME", "RS Selaras Medika"),

		WorkerConcurrency: getEnvAsInt("WORKER_CONCURRENCY", 5),
		RequestTimeout:    time.Duration(timeoutSec) * time.Second,
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
