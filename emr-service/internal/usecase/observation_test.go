package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	"github.com/azuvicenna/selaras-services/emr-service/internal/usecase"
)

type fakeObservationRepo struct {
	createErr error
	created   *domain.Observation
}

func (f *fakeObservationRepo) Create(_ context.Context, obs *domain.Observation) error {
	f.created = obs
	return f.createErr
}

func (f *fakeObservationRepo) ListByEncounter(_ context.Context, _ string, _ domain.ObservationCategory) ([]domain.Observation, error) {
	return nil, nil
}

func (f *fakeObservationRepo) ListByPatient(_ context.Context, _ domain.ObservationFilter) ([]domain.Observation, int64, error) {
	return nil, 0, nil
}

func (f *fakeObservationRepo) UpdateSatusehatID(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}

func validObservation() domain.Observation {
	temp := 37.2
	return domain.Observation{
		PatientID:      "01JQXK7M9P2R4T6V8X0Z1B3D5F",
		EncounterID:    "01ENC",
		PractitionerID: "PRAC-2026-0001",
		Category:       domain.ObservationCategoryVitalSigns,
		Code:           "8310-5",
		Name:           "Body temperature",
		ValueQuantity:  &temp,
		Unit:           "Cel",
	}
}

func TestCreateObservation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mutate  func(*domain.Observation)
		wantErr bool
	}{
		{"valid single quantity", func(*domain.Observation) {}, false},
		{"valid composite components", func(o *domain.Observation) {
			o.ValueQuantity = nil
			sys := 120.0
			o.Components = []domain.ObservationComponent{{Code: "8480-6", Name: "Systolic", ValueQuantity: &sys}}
		}, false},
		{"empty encounter id", func(o *domain.Observation) { o.EncounterID = "" }, true},
		{"empty practitioner id", func(o *domain.Observation) { o.PractitionerID = "" }, true},
		{"unspecified category", func(o *domain.Observation) { o.Category = domain.ObservationCategoryUnspecified }, true},
		{"empty name", func(o *domain.Observation) { o.Name = "" }, true},
		{"missing value and components", func(o *domain.Observation) { o.ValueQuantity = nil }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obsRepo := &fakeObservationRepo{}
			encRepo := &fakeEncounterRepo{}
			uc := usecase.NewObservationUsecase(obsRepo, encRepo)

			obs := validObservation()
			tt.mutate(&obs)

			_, err := uc.CreateObservation(ctx, obs)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Fatalf("got %v, want ErrInvalidInput", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if obsRepo.created == nil {
				t.Fatal("expected observation to be created")
			}
		})
	}
}
