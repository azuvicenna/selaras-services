package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/master-service/internal/config"
	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
	"github.com/azuvicenna/selaras-services/master-service/internal/repository"
)

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.Load().DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("database not available: %v", err)
	}
	return pool
}

func TestICD9CMRepository(t *testing.T) {
	pool := newPool(t)
	repo := repository.NewICD9CMRepository(pool)
	ctx := context.Background()

	code := fmt.Sprintf("T%d", time.Now().UnixNano()%100000)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM icd9cm WHERE code = $1", code)
	})

	item := &domain.ICD9CM{
		ID:          "01ARZ3NDEKTSV4RRFFQ69G5F01",
		Code:        code,
		Description: "Integration Test Procedure",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}

	// duplicate code error
	dup := *item
	dup.ID = "01ARZ3NDEKTSV4RRFFQ69G5F02"
	if err := repo.Create(ctx, &dup); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists on duplicate code, got %v", err)
	}

	// get by id
	got, err := repo.GetByID(ctx, item.ID)
	if err != nil || got.Code != code {
		t.Fatalf("get by id: %v, got %+v", err, got)
	}

	// get by code
	byCode, err := repo.GetByCode(ctx, code)
	if err != nil || byCode.ID != item.ID {
		t.Fatalf("get by code: %v, got %+v", err, byCode)
	}

	// list
	list, _, err := repo.List(ctx, domain.ICD9CMFilter{Limit: 1, Offset: 0})
	if err != nil || len(list) < 1 {
		t.Fatalf("list: %v, len=%d", err, len(list))
	}
}

func TestDepartmentRepository(t *testing.T) {
	pool := newPool(t)
	repo := repository.NewDepartmentRepository(pool)
	ctx := context.Background()

	code := fmt.Sprintf("DEPT_%d", time.Now().UnixNano()%100000)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM departments WHERE code = $1", code)
	})

	item := &domain.Department{
		ID:             "01ARZ3NDEKTSV4RRFFQ69G5D01",
		Code:           code,
		Name:           "Department Test",
		DepartmentType: domain.DepartmentTypeOutpatient,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.GetByCode(ctx, code)
	if err != nil || got.Name != "Department Test" {
		t.Fatalf("get by code: %v, got %+v", err, got)
	}
}
