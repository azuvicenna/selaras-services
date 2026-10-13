package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/azuvicenna/selaras-services/practitioner-service/internal/domain"
)

type mockPractitionerRepo struct {
	items map[string]*domain.Practitioner
}

func newMockPractitionerRepo() *mockPractitionerRepo {
	return &mockPractitionerRepo{items: make(map[string]*domain.Practitioner)}
}

func (m *mockPractitionerRepo) Create(ctx context.Context, item *domain.Practitioner) error {
	for _, it := range m.items {
		if it.NIK == item.NIK {
			return domain.ErrAlreadyExists
		}
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockPractitionerRepo) GetByID(ctx context.Context, id string) (*domain.Practitioner, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockPractitionerRepo) GetByNIK(ctx context.Context, nik string) (*domain.Practitioner, error) {
	for _, it := range m.items {
		if it.NIK == nik {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockPractitionerRepo) GetBySIP(ctx context.Context, sip string) (*domain.Practitioner, error) {
	for _, it := range m.items {
		if it.SIPNumber == sip {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockPractitionerRepo) List(ctx context.Context, filter domain.PractitionerFilter) ([]domain.Practitioner, int64, error) {
	var res []domain.Practitioner
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockPractitionerRepo) Update(ctx context.Context, item *domain.Practitioner) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockPractitionerRepo) UpdateStatus(ctx context.Context, id string, isActive bool) (*domain.Practitioner, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	item.IsActive = isActive
	return item, nil
}

func (m *mockPractitionerRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

func TestPractitionerUsecase_Create(t *testing.T) {
	repo := newMockPractitionerRepo()
	uc := NewPractitionerUsecase(repo)
	ctx := context.Background()

	t.Run("success create doctor", func(t *testing.T) {
		item, err := uc.CreatePractitioner(ctx, domain.Practitioner{
			NIK:            "3201012345670001",
			Name:           "dr. Budi Santoso, Sp.PD",
			Gender:         domain.GenderMale,
			Type:           domain.PractitionerTypeSpecialist,
			Specialization: "Spesialis Penyakit Dalam",
			SIPNumber:      "SIP/2026/001",
			DepartmentID:   "POLI_DALAM",
			IsActive:       true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID == "" {
			t.Errorf("expected generated ID")
		}
	})

	t.Run("success create nutritionist and physiotherapist", func(t *testing.T) {
		nutri, err := uc.CreatePractitioner(ctx, domain.Practitioner{
			NIK:            "3201011212930005",
			Name:           "Dewi Sartika, S.Gz, RD",
			Gender:         domain.GenderFemale,
			Type:           domain.PractitionerTypeNutritionist,
			Specialization: "Dietisien Klinis",
			SIPNumber:      "SIPTG/2026/001",
			DepartmentID:   "INSTALASI_GIZI",
			IsActive:       true,
		})
		if err != nil {
			t.Fatalf("unexpected error creating nutritionist: %v", err)
		}
		if nutri.Type != domain.PractitionerTypeNutritionist {
			t.Errorf("expected type Nutritionist, got %v", nutri.Type)
		}

		physio, err := uc.CreatePractitioner(ctx, domain.Practitioner{
			NIK:            "3201010304910006",
			Name:           "Bambang Irawan, S.Ft, Ftr",
			Gender:         domain.GenderMale,
			Type:           domain.PractitionerTypePhysiotherapist,
			Specialization: "Fisioterapi Muskuloskeletal",
			SIPNumber:      "SIPF/2026/001",
			DepartmentID:   "REHAB_MEDIK",
			IsActive:       true,
		})
		if err != nil {
			t.Fatalf("unexpected error creating physiotherapist: %v", err)
		}
		if physio.Type != domain.PractitionerTypePhysiotherapist {
			t.Errorf("expected type Physiotherapist, got %v", physio.Type)
		}
	})

	t.Run("missing name", func(t *testing.T) {
		_, err := uc.CreatePractitioner(ctx, domain.Practitioner{
			NIK:  "3201012345670002",
			Type: domain.PractitionerTypeDoctor,
		})
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput for missing name")
		}
	})

	t.Run("invalid nik length", func(t *testing.T) {
		_, err := uc.CreatePractitioner(ctx, domain.Practitioner{
			NIK:  "12345",
			Name: "dr. Andi",
			Type: domain.PractitionerTypeDoctor,
		})
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput for invalid nik length")
		}
	})
}

func TestPractitionerUsecase_CRUD(t *testing.T) {
	repo := newMockPractitionerRepo()
	uc := NewPractitionerUsecase(repo)
	ctx := context.Background()

	created, err := uc.CreatePractitioner(ctx, domain.Practitioner{
		NIK:          "3201012345670003",
		Name:         "dr. Siti Aminah, Sp.A",
		Type:         domain.PractitionerTypeSpecialist,
		SIPNumber:    "SIP/2026/002",
		DepartmentID: "POLI_ANAK",
		IsActive:     true,
	})
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Get by NIK
	byNIK, err := uc.GetPractitionerByNIK(ctx, "3201012345670003")
	if err != nil || byNIK.ID != created.ID {
		t.Errorf("failed to get by NIK: %v", err)
	}

	// Get by SIP
	bySIP, err := uc.GetPractitionerBySIP(ctx, "SIP/2026/002")
	if err != nil || bySIP.ID != created.ID {
		t.Errorf("failed to get by SIP: %v", err)
	}

	// Update status
	statusUpdated, err := uc.UpdatePractitionerStatus(ctx, created.ID, false)
	if err != nil || statusUpdated.IsActive != false {
		t.Errorf("failed to update status: %v", err)
	}

	// Delete
	if err := uc.DeletePractitioner(ctx, created.ID); err != nil {
		t.Errorf("failed to delete: %v", err)
	}
}
