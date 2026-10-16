package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// MigrationFile represents a parsed migration file.
type MigrationFile struct {
	Version   int64
	Name      string
	Direction string // "up" or "down"
	FilePath  string
}

// ParseMigrationFileName parses migration files in format: {version}_{name}.{direction}.sql
func ParseMigrationFileName(filename string) (*MigrationFile, error) {
	cleanName := filepath.Base(filename)
	parts := strings.Split(cleanName, ".")
	if len(parts) != 3 || parts[2] != "sql" {
		return nil, fmt.Errorf("invalid migration filename: %s (expected {version}_{name}.{up|down}.sql)", cleanName)
	}

	direction := strings.ToLower(parts[1])
	if direction != "up" && direction != "down" {
		return nil, fmt.Errorf("invalid migration direction: %s (must be up or down)", direction)
	}

	prefix := parts[0]
	underscoreIdx := strings.Index(prefix, "_")
	if underscoreIdx == -1 {
		return nil, fmt.Errorf("invalid migration prefix: %s (missing underscore separator)", prefix)
	}

	versionStr := prefix[:underscoreIdx]
	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid version %q: %w", versionStr, err)
	}

	name := prefix[underscoreIdx+1:]

	return &MigrationFile{
		Version:   version,
		Name:      name,
		Direction: direction,
		FilePath:  filename,
	}, nil
}

// EnsureSchemaTable creates the tracking table in PostgreSQL if it does not exist.
func EnsureSchemaTable(ctx context.Context, conn *pgx.Conn) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version BIGINT PRIMARY KEY,
		name VARCHAR(255) NOT NULL DEFAULT '',
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);`
	_, err := conn.Exec(ctx, query)
	return err
}

// GetAppliedVersions loads all applied migration versions from schema_migrations.
func GetAppliedVersions(ctx context.Context, conn *pgx.Conn) (map[int64]time.Time, error) {
	rows, err := conn.Query(ctx, "SELECT version, applied_at FROM schema_migrations ORDER BY version ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int64]time.Time)
	for rows.Next() {
		var v int64
		var t time.Time
		if err := rows.Scan(&v, &t); err != nil {
			return nil, err
		}
		applied[v] = t
	}
	return applied, rows.Err()
}

// LoadMigrationFiles scans and parses all sql files in the migrations directory.
func LoadMigrationFiles(dir, direction string) ([]MigrationFile, error) {
	pattern := filepath.Join(dir, fmt.Sprintf("*.%s.sql", direction))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var files []MigrationFile
	for _, m := range matches {
		mf, err := ParseMigrationFileName(m)
		if err != nil {
			return nil, err
		}
		if mf.Direction == direction {
			files = append(files, *mf)
		}
	}

	// Sort files by version
	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})

	return files, nil
}

// StatusResult stores the migration status of a service.
type StatusResult struct {
	Service string
	Applied int
	Pending int
	Total   int
	Latest  int64
}

// MigrateServiceUp executes all unapplied .up.sql files for a service.
func MigrateServiceUp(ctx context.Context, svc ServiceConfig, host string, port int, ssl string) error {
	dbURL := svc.DatabaseURL(host, port, ssl)
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", svc.DBName, err)
	}
	defer conn.Close(ctx)

	if err := EnsureSchemaTable(ctx, conn); err != nil {
		return fmt.Errorf("failed to ensure schema_migrations in %s: %w", svc.DBName, err)
	}

	applied, err := GetAppliedVersions(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to get applied versions for %s: %w", svc.DBName, err)
	}

	files, err := LoadMigrationFiles(svc.MigrationsDir, "up")
	if err != nil {
		return fmt.Errorf("failed to load migration files for %s: %w", svc.Name, err)
	}

	if len(files) == 0 {
		slog.Warn("No migrations found for service", slog.String("service", svc.Name), slog.String("dir", svc.MigrationsDir))
		return nil
	}

	appliedCount := 0
	for _, mf := range files {
		if _, exists := applied[mf.Version]; exists {
			continue
		}

		sqlBytes, err := os.ReadFile(mf.FilePath)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", mf.FilePath, err)
		}

		slog.Info("Applying migration",
			slog.String("service", svc.Name),
			slog.Int64("version", mf.Version),
			slog.String("name", mf.Name),
		)

		tx, err := conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("execute migration %s: %w", filepath.Base(mf.FilePath), err)
		}

		_, err = tx.Exec(ctx, "INSERT INTO schema_migrations (version, name, applied_at) VALUES ($1, $2, now())",
			mf.Version, mf.Name)
		if err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("record migration %d: %w", mf.Version, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %d: %w", mf.Version, err)
		}

		appliedCount++
	}

	if appliedCount == 0 {
		slog.Info("Database schema is already up to date", slog.String("service", svc.Name))
	} else {
		slog.Info("Successfully applied migrations", slog.String("service", svc.Name), slog.Int("applied", appliedCount))
	}
	return nil
}

// MigrateServiceDown rolls back the most recent applied migration for a service.
func MigrateServiceDown(ctx context.Context, svc ServiceConfig, host string, port int, ssl string) error {
	dbURL := svc.DatabaseURL(host, port, ssl)
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", svc.DBName, err)
	}
	defer conn.Close(ctx)

	if err := EnsureSchemaTable(ctx, conn); err != nil {
		return err
	}

	applied, err := GetAppliedVersions(ctx, conn)
	if err != nil {
		return err
	}

	if len(applied) == 0 {
		slog.Info("No migrations to rollback", slog.String("service", svc.Name))
		return nil
	}

	// Find maximum applied version
	var maxVersion int64
	for v := range applied {
		if v > maxVersion {
			maxVersion = v
		}
	}

	files, err := LoadMigrationFiles(svc.MigrationsDir, "down")
	if err != nil {
		return err
	}

	var targetFile *MigrationFile
	for i := range files {
		if files[i].Version == maxVersion {
			targetFile = &files[i]
			break
		}
	}

	if targetFile == nil {
		return fmt.Errorf("corresponding .down.sql for version %d not found in %s", maxVersion, svc.MigrationsDir)
	}

	sqlBytes, err := os.ReadFile(targetFile.FilePath)
	if err != nil {
		return fmt.Errorf("read file %s: %w", targetFile.FilePath, err)
	}

	slog.Info("Rolling back migration",
		slog.String("service", svc.Name),
		slog.Int64("version", targetFile.Version),
		slog.String("name", targetFile.Name),
	)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("execute rollback %s: %w", filepath.Base(targetFile.FilePath), err)
	}

	if _, err := tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", maxVersion); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("delete migration record %d: %w", maxVersion, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	slog.Info("Successfully rolled back migration",
		slog.String("service", svc.Name),
		slog.Int64("version", maxVersion),
	)
	return nil
}

// GetServiceStatus checks applied vs pending migration counts.
func GetServiceStatus(ctx context.Context, svc ServiceConfig, host string, port int, ssl string) (*StatusResult, error) {
	dbURL := svc.DatabaseURL(host, port, ssl)
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", svc.DBName, err)
	}
	defer conn.Close(ctx)

	if err := EnsureSchemaTable(ctx, conn); err != nil {
		return nil, err
	}

	applied, err := GetAppliedVersions(ctx, conn)
	if err != nil {
		return nil, err
	}

	files, err := LoadMigrationFiles(svc.MigrationsDir, "up")
	if err != nil {
		return nil, err
	}

	var latestVersion int64
	for v := range applied {
		if v > latestVersion {
			latestVersion = v
		}
	}

	pendingCount := 0
	for _, f := range files {
		if _, ok := applied[f.Version]; !ok {
			pendingCount++
		}
	}

	return &StatusResult{
		Service: svc.Name,
		Applied: len(applied),
		Pending: pendingCount,
		Total:   len(files),
		Latest:  latestVersion,
	}, nil
}
