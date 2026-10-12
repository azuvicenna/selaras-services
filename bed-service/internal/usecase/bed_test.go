package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

// --- Mock WardClass Repository ---
type mockWardClassRepo struct {
	items map[string]*domain.WardClass
}

func newMockWardClassRepo() *mockWardClassRepo {
	return &mockWardClassRepo{items: make(map[string]*domain.WardClass)}
}

func (m *mockWardClassRepo) Create(ctx context.Context, item *domain.WardClass) error {
	for _, it := range m.items {
		if it.Code == item.Code {
			return domain.ErrAlreadyExists
		}
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockWardClassRepo) GetByID(ctx context.Context, id string) (*domain.WardClass, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockWardClassRepo) GetByCode(ctx context.Context, code string) (*domain.WardClass, error) {
	for _, it := range m.items {
		if it.Code == code {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockWardClassRepo) List(ctx context.Context, filter domain.WardClassFilter) ([]domain.WardClass, int64, error) {
	var res []domain.WardClass
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockWardClassRepo) Update(ctx context.Context, item *domain.WardClass) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockWardClassRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

// --- Mock Ward Repository ---
type mockWardRepo struct {
	items map[string]*domain.Ward
}

func newMockWardRepo() *mockWardRepo {
	return &mockWardRepo{items: make(map[string]*domain.Ward)}
}

func (m *mockWardRepo) Create(ctx context.Context, item *domain.Ward) error {
	for _, it := range m.items {
		if it.Code == item.Code {
			return domain.ErrAlreadyExists
		}
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockWardRepo) GetByID(ctx context.Context, id string) (*domain.Ward, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockWardRepo) GetByCode(ctx context.Context, code string) (*domain.Ward, error) {
	for _, it := range m.items {
		if it.Code == code {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockWardRepo) List(ctx context.Context, filter domain.WardFilter) ([]domain.Ward, int64, error) {
	var res []domain.Ward
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockWardRepo) Update(ctx context.Context, item *domain.Ward) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockWardRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

// --- Mock Bed Repository ---
type mockBedRepo struct {
	items map[string]*domain.Bed
}

func newMockBedRepo() *mockBedRepo {
	return &mockBedRepo{items: make(map[string]*domain.Bed)}
}

func (m *mockBedRepo) Create(ctx context.Context, item *domain.Bed) error {
	for _, it := range m.items {
		if it.WardID == item.WardID && it.Code == item.Code {
			return domain.ErrAlreadyExists
		}
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockBedRepo) GetByID(ctx context.Context, id string) (*domain.Bed, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockBedRepo) GetByCode(ctx context.Context, wardID, code string) (*domain.Bed, error) {
	for _, it := range m.items {
		if it.WardID == wardID && it.Code == code {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockBedRepo) List(ctx context.Context, filter domain.BedFilter) ([]domain.Bed, int64, error) {
	var res []domain.Bed
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockBedRepo) Update(ctx context.Context, item *domain.Bed) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockBedRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

func (m *mockBedRepo) UpdateStatus(ctx context.Context, id string, status domain.BedStatus, updatedAt time.Time) error {
	item, ok := m.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	item.Status = status
	item.UpdatedAt = updatedAt
	return nil
}

func (m *mockBedRepo) GetSummary(ctx context.Context, wardClassID string) (*domain.BedSummary, error) {
	summary := &domain.BedSummary{}
	for _, it := range m.items {
		summary.TotalBeds++
		switch it.Status {
		case domain.BedStatusAvailable:
			summary.AvailableBeds++
		case domain.BedStatusOccupied:
			summary.OccupiedBeds++
		}
	}
	return summary, nil
}

// -----------------------------------------------------------------------------
// TESTS
// -----------------------------------------------------------------------------

func TestWardClassUsecase(t *testing.T) {
	repo := newMockWardClassRepo()
	uc := NewWardClassUsecase(repo)
	ctx := context.Background()

	// 1. Validation fail
	_, err := uc.CreateWardClass(ctx, domain.WardClass{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	// 2. Success with ULID generation
	created, err := uc.CreateWardClass(ctx, domain.WardClass{
		Code:                "VIP",
		Name:                "VIP Standard",
		BPJSEquivalentClass: "VIP",
		IsActive:            true,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if len(created.ID) != 26 {
		t.Fatalf("expected ULID length 26, got %d (id=%s)", len(created.ID), created.ID)
	}

	// 3. Get by ID
	got, err := uc.GetWardClassByID(ctx, created.ID)
	if err != nil || got.Code != "VIP" {
		t.Fatalf("get by id failed: %v", err)
	}

	// 4. Duplicate code error
	_, err = uc.CreateWardClass(ctx, domain.WardClass{
		Code: "VIP",
		Name: "Duplicate VIP",
	})
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}

	// 5. Update
	created.Name = "VIP Deluxe"
	updated, err := uc.UpdateWardClass(ctx, *created)
	if err != nil || updated.Name != "VIP Deluxe" {
		t.Fatalf("update failed: %v", err)
	}

	// 6. Delete
	if err := uc.DeleteWardClass(ctx, created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	_, err = uc.GetWardClassByID(ctx, created.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestWardUsecase(t *testing.T) {
	repo := newMockWardRepo()
	uc := NewWardUsecase(repo)
	ctx := context.Background()

	// 1. Validation fail
	_, err := uc.CreateWard(ctx, domain.Ward{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	// 2. Success with defaults
	created, err := uc.CreateWard(ctx, domain.Ward{
		Code:        "MAWAR-01",
		Name:        "Kamar Mawar 101",
		WardClassID: "01ARZ3NDEKTSV4RRFFQ69G5F01",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if len(created.ID) != 26 {
		t.Fatalf("expected ULID length 26, got %d", len(created.ID))
	}
	if created.Floor != 1 || created.Capacity != 1 {
		t.Fatalf("expected floor=1 and capacity=1 defaults, got floor=%d, cap=%d", created.Floor, created.Capacity)
	}

	// 3. Get by code
	byCode, err := uc.GetWardByCode(ctx, "MAWAR-01")
	if err != nil || byCode.Name != "Kamar Mawar 101" {
		t.Fatalf("get by code failed: %v", err)
	}
}

func TestBedUsecase(t *testing.T) {
	repo := newMockBedRepo()
	uc := NewBedUsecase(repo)
	ctx := context.Background()

	// 1. Validation fail
	_, err := uc.CreateBed(ctx, domain.Bed{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	// 2. Success with default status Available
	created, err := uc.CreateBed(ctx, domain.Bed{
		Code:   "MWR-101-A",
		WardID: "01ARZ3NDEKTSV4RRFFQ69G5W01",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if len(created.ID) != 26 {
		t.Fatalf("expected ULID length 26, got %d", len(created.ID))
	}
	if created.Status != domain.BedStatusAvailable {
		t.Fatalf("expected default status Available (1), got %v", created.Status)
	}

	// 3. Update status
	if err := uc.UpdateBedStatus(ctx, created.ID, domain.BedStatusOccupied); err != nil {
		t.Fatalf("update status failed: %v", err)
	}
	got, err := uc.GetBedByID(ctx, created.ID)
	if err != nil || got.Status != domain.BedStatusOccupied {
		t.Fatalf("status update did not reflect, got status=%v", got.Status)
	}

	// 4. Invalid status update
	if err := uc.UpdateBedStatus(ctx, created.ID, 99); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for invalid status, got %v", err)
	}

	// 5. Summary
	summary, err := uc.GetBedSummary(ctx, "")
	if err != nil || summary.TotalBeds != 1 || summary.OccupiedBeds != 1 {
		t.Fatalf("summary failed: %+v", summary)
	}
}
