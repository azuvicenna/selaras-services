package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                string
	MetricsPort         string
	Env                 string
	CORSAllowedOrigins  string
	JWTSecret           string
	APIKeySecret        string
	RateLimitRPS        int
	RateLimitBurst      int
	EnableSignatureCheck bool

	// Downstream gRPC Services
	AuthServiceAddr         string
	PatientServiceAddr      string
	PractitionerServiceAddr string
	AdmissionServiceAddr    string
	BedServiceAddr          string
	MasterServiceAddr       string
	EMRServiceAddr          string
	BillingServiceAddr      string
	PharmacyServiceAddr     string
	LaboratoryServiceAddr   string
	RadiologyServiceAddr    string
	DocumentServiceAddr     string
}

func Load() *Config {
	return &Config{
		Port:                 getEnv("HTTP_PORT", "8082"),
		MetricsPort:          getEnv("METRICS_PORT", "9092"),
		Env:                  getEnv("ENV", "development"),
		CORSAllowedOrigins:   getEnv("CORS_ALLOWED_ORIGINS", "*"),
		JWTSecret:            getEnv("JWT_SECRET", "selaras-hospital-jwt-secret-key-32bytes!"),
		APIKeySecret:         getEnv("B2B_API_KEY_SECRET", "selaras-b2b-secret-salt-key-32bytes!"),
		RateLimitRPS:         getEnvInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst:       getEnvInt("RATE_LIMIT_BURST", 200),
		EnableSignatureCheck: getEnvBool("ENABLE_SIGNATURE_CHECK", false),

		AuthServiceAddr:         getEnv("AUTH_SERVICE_ADDR", "localhost:50057"),
		PatientServiceAddr:      getEnv("PATIENT_SERVICE_ADDR", "localhost:50051"),
		PractitionerServiceAddr: getEnv("PRACTITIONER_SERVICE_ADDR", "localhost:50056"),
		AdmissionServiceAddr:    getEnv("ADMISSION_SERVICE_ADDR", "localhost:50055"),
		BedServiceAddr:          getEnv("BED_SERVICE_ADDR", "localhost:50054"),
		MasterServiceAddr:       getEnv("MASTER_SERVICE_ADDR", "localhost:50053"),
		EMRServiceAddr:          getEnv("EMR_SERVICE_ADDR", "localhost:50052"),
		BillingServiceAddr:      getEnv("BILLING_SERVICE_ADDR", "localhost:50058"),
		PharmacyServiceAddr:     getEnv("PHARMACY_SERVICE_ADDR", "localhost:50059"),
		LaboratoryServiceAddr:   getEnv("LABORATORY_SERVICE_ADDR", "localhost:50060"),
		RadiologyServiceAddr:    getEnv("RADIOLOGY_SERVICE_ADDR", "localhost:50061"),
		DocumentServiceAddr:     getEnv("DOCUMENT_SERVICE_ADDR", "localhost:50062"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return fallback
}
