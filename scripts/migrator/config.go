package main

import (
	"fmt"
	"path/filepath"
)

// ServiceConfig defines the database credentials, directories, and gRPC endpoints for a microservice.
type ServiceConfig struct {
	Name          string
	DBName        string
	DBUser        string
	DBPassword    string
	GRPCPort      int
	MigrationsDir string
	SeedersDir    string
	SeedFile      string
}

// DatabaseURL generates the PostgreSQL connection string for this service.
func (s ServiceConfig) DatabaseURL(host string, port int, sslMode string) string {
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		s.DBUser, s.DBPassword, host, port, s.DBName, sslMode)
}

// DefaultServices defines all 12 core hospital services in strict dependency order.
var DefaultServices = []ServiceConfig{
	{
		Name:          "auth-service",
		DBName:        "auth_db",
		DBUser:        "auth",
		DBPassword:    "change-me",
		GRPCPort:      50057,
		MigrationsDir: filepath.Join("auth-service", "db", "migrations"),
		SeedersDir:    filepath.Join("auth-service", "db", "seeders"),
		SeedFile:      "users.json",
	},
	{
		Name:          "master-service",
		DBName:        "master_db",
		DBUser:        "master",
		DBPassword:    "change-me",
		GRPCPort:      50053,
		MigrationsDir: filepath.Join("master-service", "db", "migrations"),
		SeedersDir:    filepath.Join("master-service", "db", "seeders"),
		SeedFile:      "master_data.json",
	},
	{
		Name:          "practitioner-service",
		DBName:        "practitioner_db",
		DBUser:        "practitioner",
		DBPassword:    "change-me",
		GRPCPort:      50056,
		MigrationsDir: filepath.Join("practitioner-service", "db", "migrations"),
		SeedersDir:    filepath.Join("practitioner-service", "db", "seeders"),
		SeedFile:      "practitioners.json",
	},
	{
		Name:          "patient-service",
		DBName:        "patient_db",
		DBUser:        "patient",
		DBPassword:    "change-me",
		GRPCPort:      50051,
		MigrationsDir: filepath.Join("patient-service", "db", "migrations"),
		SeedersDir:    filepath.Join("patient-service", "db", "seeders"),
		SeedFile:      "patients.json",
	},
	{
		Name:          "bed-service",
		DBName:        "bed_db",
		DBUser:        "bed",
		DBPassword:    "change-me",
		GRPCPort:      50054,
		MigrationsDir: filepath.Join("bed-service", "db", "migrations"),
		SeedersDir:    filepath.Join("bed-service", "db", "seeders"),
		SeedFile:      "bed_data.json",
	},
	{
		Name:          "admission-service",
		DBName:        "admission_db",
		DBUser:        "admission",
		DBPassword:    "change-me",
		GRPCPort:      50055,
		MigrationsDir: filepath.Join("admission-service", "db", "migrations"),
		SeedersDir:    filepath.Join("admission-service", "db", "seeders"),
		SeedFile:      "admission_data.json",
	},
	{
		Name:          "emr-service",
		DBName:        "emr_db",
		DBUser:        "emr",
		DBPassword:    "change-me",
		GRPCPort:      50052,
		MigrationsDir: filepath.Join("emr-service", "db", "migrations"),
		SeedersDir:    filepath.Join("emr-service", "db", "seeders"),
		SeedFile:      "emr.json",
	},
	{
		Name:          "billing-service",
		DBName:        "billing_db",
		DBUser:        "billing",
		DBPassword:    "change-me",
		GRPCPort:      50058,
		MigrationsDir: filepath.Join("billing-service", "db", "migrations"),
		SeedersDir:    filepath.Join("billing-service", "db", "seeders"),
		SeedFile:      "tariffs.json",
	},
	{
		Name:          "pharmacy-service",
		DBName:        "pharmacy_db",
		DBUser:        "pharmacy",
		DBPassword:    "change-me",
		GRPCPort:      50059,
		MigrationsDir: filepath.Join("pharmacy-service", "db", "migrations"),
		SeedersDir:    filepath.Join("pharmacy-service", "db", "seeders"),
		SeedFile:      "medications.json",
	},
	{
		Name:          "laboratory-service",
		DBName:        "laboratory_db",
		DBUser:        "laboratory",
		DBPassword:    "change-me",
		GRPCPort:      50060,
		MigrationsDir: filepath.Join("laboratory-service", "db", "migrations"),
		SeedersDir:    filepath.Join("laboratory-service", "db", "seeders"),
		SeedFile:      "tests.json",
	},
	{
		Name:          "radiology-service",
		DBName:        "radiology_db",
		DBUser:        "radiology",
		DBPassword:    "change-me",
		GRPCPort:      50061,
		MigrationsDir: filepath.Join("radiology-service", "db", "migrations"),
		SeedersDir:    filepath.Join("radiology-service", "db", "seeders"),
		SeedFile:      "exams.json",
	},
	{
		Name:          "document-service",
		DBName:        "document_db",
		DBUser:        "document",
		DBPassword:    "change-me",
		GRPCPort:      50062,
		MigrationsDir: filepath.Join("document-service", "db", "migrations"),
		SeedersDir:    filepath.Join("document-service", "db", "seeders"),
		SeedFile:      "templates.json",
	},
}

// FindService returns the configuration for a given service name.
func FindService(name string) (*ServiceConfig, error) {
	for i := range DefaultServices {
		if DefaultServices[i].Name == name {
			return &DefaultServices[i], nil
		}
	}
	return nil, fmt.Errorf("service %q not found in registry", name)
}
