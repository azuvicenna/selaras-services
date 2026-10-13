package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azuvicenna/selaras-services/practitioner-service/internal/config"
	"github.com/azuvicenna/selaras-services/practitioner-service/internal/domain"
	"github.com/azuvicenna/selaras-services/practitioner-service/internal/repository"
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

func TestPractitionerRepository(t *testing.T) {
	pool := newPool(t)
	repo := repository.NewPractitionerRepository(pool)
	ctx := context.Background()

	id := fmt.Sprintf("01ARZ%021d", time.Now().UnixNano()%1000000)
	nik := fmt.Sprintf("320101%010d", time.Now().UnixNano()%10000000000)
	sip := fmt.Sprintf("SIP/TEST/%d", time.Now().UnixNano()%100000)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM practitioners WHERE id = $1", id)
	})

	item := &domain.Practitioner{
		ID:             id,
		NIK:            nik,
		Name:           "dr. Test Practitioner, Sp.PD",
		Gender:         domain.GenderMale,
		BirthDate:      "1985-05-15",
		Phone:          "+6281234567890",
		Email:          "test@selaras.id",
		Type:           domain.PractitionerTypeSpecialist,
		Specialization: "Penyakit Dalam",
		SIPNumber:      sip,
		SIPExpiryDate:  "2029-05-15",
		STRNumber:      "31.1.1.100.1.20.999999",
		STRExpiryDate:  "2029-05-15",
		DepartmentID:   "POLI_DALAM",
		IHSNumber:      "N10009999",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}

	// Test GetByID
	found, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if found.Name != item.Name {
		t.Errorf("expected name %s, got %s", item.Name, found.Name)
	}

	// Test GetByNIK
	byNIK, err := repo.GetByNIK(ctx, nik)
	if err != nil || byNIK.ID != id {
		t.Errorf("get by NIK failed: %v", err)
	}

	// Test GetBySIP
	bySIP, err := repo.GetBySIP(ctx, sip)
	if err != nil || bySIP.ID != id {
		t.Errorf("get by SIP failed: %v", err)
	}

	// Test unique NIK error
	dup := *item
	dup.ID = fmt.Sprintf("01ARZ%021d", (time.Now().UnixNano()+1)%1000000)
	if err := repo.Create(ctx, &dup); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists on duplicate NIK, got: %v", err)
	}

	// Test Delete
	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
