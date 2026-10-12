package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/bed-service/internal/config"
	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
	"github.com/azuvicenna/selaras-services/bed-service/internal/repository"
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

func TestWardClassRepository(t *testing.T) {
	pool := newPool(t)
	repo := repository.NewWardClassRepository(pool)
	ctx := context.Background()

	code := fmt.Sprintf("TC%d", time.Now().UnixNano()%100000)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM ward_classes WHERE code = $1", code)
	})

	item := &domain.WardClass{
		ID:                  "01ARZ3NDEKTSV4RRFFQ69G5C01",
		Code:                code,
		Name:                "Test Ward Class",
		BPJSEquivalentClass: "VIP",
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}

	// duplicate code error
	dup := *item
	dup.ID = "01ARZ3NDEKTSV4RRFFQ69G5C02"
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
	list, _, err := repo.List(ctx, domain.WardClassFilter{Limit: 1, Offset: 0})
	if err != nil || len(list) < 1 {
		t.Fatalf("list: %v, len=%d", err, len(list))
	}
}
