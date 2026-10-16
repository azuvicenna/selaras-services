package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMigrationFileName(t *testing.T) {
	tests := []struct {
		filename      string
		wantVersion   int64
		wantName      string
		wantDirection string
		wantErr       bool
	}{
		{
			filename:      "000001_create_auth_tables.up.sql",
			wantVersion:   1,
			wantName:      "create_auth_tables",
			wantDirection: "up",
			wantErr:       false,
		},
		{
			filename:      "000008_create_diagnostic_reports_table.down.sql",
			wantVersion:   8,
			wantName:      "create_diagnostic_reports_table",
			wantDirection: "down",
			wantErr:       false,
		},
		{
			filename: "invalid_file.txt",
			wantErr:  true,
		},
		{
			filename: "abc_test.up.sql",
			wantErr:  true,
		},
		{
			filename: "000001_no_direction.sql",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			res, err := ParseMigrationFileName(tt.filename)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %s, got nil", tt.filename)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.Version != tt.wantVersion {
				t.Errorf("got version %d, want %d", res.Version, tt.wantVersion)
			}
			if res.Name != tt.wantName {
				t.Errorf("got name %q, want %q", res.Name, tt.wantName)
			}
			if res.Direction != tt.wantDirection {
				t.Errorf("got direction %q, want %q", res.Direction, tt.wantDirection)
			}
		})
	}
}

func TestDefaultServicesConfiguration(t *testing.T) {
	if len(DefaultServices) != 12 {
		t.Fatalf("expected exactly 12 core services, got %d", len(DefaultServices))
	}

	// Change working directory reference relative to workspace root if running from package dir
	rootDir := filepath.Join("..", "..")
	if _, err := os.Stat("admission-service"); err == nil {
		rootDir = "."
	}

	for _, svc := range DefaultServices {
		t.Run(svc.Name, func(t *testing.T) {
			if svc.DBName == "" {
				t.Errorf("service %s has empty DBName", svc.Name)
			}
			if svc.DBUser == "" {
				t.Errorf("service %s has empty DBUser", svc.Name)
			}
			if svc.DBPassword == "" {
				t.Errorf("service %s has empty DBPassword", svc.Name)
			}
			if svc.GRPCPort <= 0 {
				t.Errorf("service %s has invalid GRPCPort %d", svc.Name, svc.GRPCPort)
			}

			// Verify migration directory exists on disk
			dir := filepath.Join(rootDir, svc.MigrationsDir)
			info, err := os.Stat(dir)
			if err != nil || !info.IsDir() {
				t.Errorf("service %s migrations dir %s does not exist", svc.Name, dir)
			}

			// Verify at least one migration file exists
			files, err := LoadMigrationFiles(dir, "up")
			if err != nil {
				t.Errorf("failed to load migration files for %s: %v", svc.Name, err)
			}
			if len(files) == 0 {
				t.Errorf("service %s has 0 up migrations", svc.Name)
			}
		})
	}
}

func TestFindService(t *testing.T) {
	svc, err := FindService("auth-service")
	if err != nil || svc == nil || svc.DBName != "auth_db" {
		t.Fatalf("failed to find auth-service: %v", err)
	}

	_, err = FindService("non-existent-service")
	if err == nil {
		t.Fatalf("expected error for non-existent service")
	}
}
