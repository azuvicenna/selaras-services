package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

var nikPattern = regexp.MustCompile(`^\d{16}$`)

type PatientUsecase struct {
	repo domain.PatientRepository
}

func NewPatientUsecase(repo domain.PatientRepository) *PatientUsecase {
	return &PatientUsecase{repo: repo}
}

func (u *PatientUsecase) Register(ctx context.Context, p domain.Patient) (*domain.Patient, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	p.ID = ulid.Make().String()
	if err := u.repo.Create(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (u *PatientUsecase) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *PatientUsecase) List(ctx context.Context, limit, offset int) ([]domain.Patient, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	return u.repo.List(ctx, min(limit, maxLimit), max(offset, 0))
}

func (u *PatientUsecase) Update(ctx context.Context, p domain.Patient) (*domain.Patient, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	if err := u.repo.Update(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func validate(p domain.Patient) error {
	switch {
	case strings.TrimSpace(p.Name) == "":
		return fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	case !nikPattern.MatchString(p.NIK):
		return fmt.Errorf("%w: nik must be 16 digits", domain.ErrInvalidInput)
	case p.BirthDate.IsZero() || p.BirthDate.After(time.Now()):
		return fmt.Errorf("%w: birth date is invalid", domain.ErrInvalidInput)
	case p.Gender != domain.Male && p.Gender != domain.Female:
		return fmt.Errorf("%w: gender must be male or female", domain.ErrInvalidInput)
	}
	return nil
}