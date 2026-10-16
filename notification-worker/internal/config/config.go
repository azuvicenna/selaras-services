package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPPort    string
	MetricsPort string

	// Apache Kafka Settings
	KafkaBrokers       []string
	KafkaTopic         string
	KafkaConsumerGroup string
	KafkaDLQTopic      string
	KafkaMaxRetries    int
	KafkaMinBytes      int
	KafkaMaxBytes      int
	KafkaMaxWait       time.Duration

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
	brokersStr := getEnv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(brokersStr, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	timeoutSec := getEnvAsInt("REQUEST_TIMEOUT_SEC", 30)

	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8087"),
		MetricsPort:        getEnv("METRICS_PORT", "9097"),
		KafkaBrokers:       brokers,
		KafkaTopic:         getEnv("KAFKA_TOPIC", "hospital.notifications"),
		KafkaConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "notification-worker-group"),
		KafkaDLQTopic:      getEnv("KAFKA_DLQ_TOPIC", "hospital.notifications.dlq"),
		KafkaMaxRetries:    getEnvAsInt("KAFKA_MAX_RETRIES", 3),
		KafkaMinBytes:      getEnvAsInt("KAFKA_MIN_BYTES", 10e3), // 10KB
		KafkaMaxBytes:      getEnvAsInt("KAFKA_MAX_BYTES", 10e6), // 10MB
		KafkaMaxWait:       time.Duration(getEnvAsInt("KAFKA_MAX_WAIT_MS", 1000)) * time.Millisecond,

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
