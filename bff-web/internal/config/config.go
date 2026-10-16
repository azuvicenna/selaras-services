package config

import (
	"os"
)

type Config struct {
	HTTPPort           string
	MetricsPort        string
	CORSAllowedOrigins string
	JWTSecret          string

	// gRPC Service Endpoints
	AuthServiceAddr         string
	PatientServiceAddr      string
	PractitionerServiceAddr string
	AdmissionServiceAddr    string
	BedServiceAddr          string
	MasterServiceAddr       string
	EmrServiceAddr          string
	BillingServiceAddr      string
	PharmacyServiceAddr     string
	LaboratoryServiceAddr   string
	RadiologyServiceAddr    string
	DocumentServiceAddr     string
}

func Load() *Config {
	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		MetricsPort:        getEnv("METRICS_PORT", "9090"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		JWTSecret:          getEnv("JWT_SECRET", "selaras-hospital-jwt-secret-key-32bytes!"),

		AuthServiceAddr:         getEnv("AUTH_SERVICE_ADDR", "localhost:50057"),
		PatientServiceAddr:      getEnv("PATIENT_SERVICE_ADDR", "localhost:50051"),
		PractitionerServiceAddr: getEnv("PRACTITIONER_SERVICE_ADDR", "localhost:50056"),
		AdmissionServiceAddr:    getEnv("ADMISSION_SERVICE_ADDR", "localhost:50055"),
		BedServiceAddr:          getEnv("BED_SERVICE_ADDR", "localhost:50054"),
		MasterServiceAddr:       getEnv("MASTER_SERVICE_ADDR", "localhost:50053"),
		EmrServiceAddr:          getEnv("EMR_SERVICE_ADDR", "localhost:50052"),
		BillingServiceAddr:      getEnv("BILLING_SERVICE_ADDR", "localhost:50058"),
		PharmacyServiceAddr:     getEnv("PHARMACY_SERVICE_ADDR", "localhost:50059"),
		LaboratoryServiceAddr:   getEnv("LABORATORY_SERVICE_ADDR", "localhost:50060"),
		RadiologyServiceAddr:    getEnv("RADIOLOGY_SERVICE_ADDR", "localhost:50061"),
		DocumentServiceAddr:     getEnv("DOCUMENT_SERVICE_ADDR", "localhost:50062"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
