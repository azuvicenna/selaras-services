package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
	"github.com/azuvicenna/selaras-services/bed-service/internal/event"
)

var _ domain.BedUsecase = (*BedUsecaseImpl)(nil)

type BedUsecaseImpl struct {
	repo      domain.BedRepository
	publisher event.EventPublisher
}

func NewBedUsecase(repo domain.BedRepository) *BedUsecaseImpl {
	return &BedUsecaseImpl{
		repo:      repo,
		publisher: event.NewNoopPublisher(),
	}
}

func (u *BedUsecaseImpl) WithEventPublisher(p event.EventPublisher) *BedUsecaseImpl {
	if p != nil {
		u.publisher = p
	}
	return u
}

func (u *BedUsecaseImpl) CreateBed(ctx context.Context, item domain.Bed) (*domain.Bed, error) {
	item.Code = strings.TrimSpace(item.Code)
	item.WardID = strings.TrimSpace(item.WardID)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: bed code is required", domain.ErrInvalidInput)
	}
	if item.WardID == "" {
		return nil, fmt.Errorf("%w: ward id is required", domain.ErrInvalidInput)
	}

	if item.Status == domain.BedStatusUnspecified {
		item.Status = domain.BedStatusAvailable
	}

	if item.ID == "" {
		item.ID = ulid.Make().String()
	}
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := u.repo.Create(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *BedUsecaseImpl) GetBedByID(ctx context.Context, id string) (*domain.Bed, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByID(ctx, id)
}

func (u *BedUsecaseImpl) GetBedByCode(ctx context.Context, wardID, code string) (*domain.Bed, error) {
	wardID = strings.TrimSpace(wardID)
	code = strings.TrimSpace(code)
	if wardID == "" {
		return nil, fmt.Errorf("%w: ward id is required", domain.ErrInvalidInput)
	}
	if code == "" {
		return nil, fmt.Errorf("%w: bed code is required", domain.ErrInvalidInput)
	}
	return u.repo.GetByCode(ctx, wardID, code)
}

func (u *BedUsecaseImpl) ListBeds(ctx context.Context, filter domain.BedFilter) ([]domain.Bed, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	filter.Limit = min(filter.Limit, maxLimit)
	filter.Offset = max(filter.Offset, 0)
	filter.Query = strings.TrimSpace(filter.Query)
	filter.WardID = strings.TrimSpace(filter.WardID)

	return u.repo.List(ctx, filter)
}

func (u *BedUsecaseImpl) UpdateBed(ctx context.Context, item domain.Bed) (*domain.Bed, error) {
	if strings.TrimSpace(item.ID) == "" {
		return nil, fmt.Errorf("%w: id is required for update", domain.ErrInvalidInput)
	}
	item.Code = strings.TrimSpace(item.Code)
	item.WardID = strings.TrimSpace(item.WardID)

	if item.Code == "" {
		return nil, fmt.Errorf("%w: bed code is required", domain.ErrInvalidInput)
	}
	if item.WardID == "" {
		return nil, fmt.Errorf("%w: ward id is required", domain.ErrInvalidInput)
	}
	if item.Status == domain.BedStatusUnspecified {
		item.Status = domain.BedStatusAvailable
	}

	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (u *BedUsecaseImpl) DeleteBed(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	return u.repo.Delete(ctx, id)
}

func (u *BedUsecaseImpl) UpdateBedStatus(ctx context.Context, id string, status domain.BedStatus) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if status <= domain.BedStatusUnspecified || status > domain.BedStatusIsolation {
		return fmt.Errorf("%w: invalid bed status", domain.ErrInvalidInput)
	}

	now := time.Now()
	if err := u.repo.UpdateStatus(ctx, id, status, now); err != nil {
		return err
	}

	if u.publisher != nil {
		pubCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		tersedia := 0
		if status == domain.BedStatusAvailable {
			tersedia = 1
		}

		aplicaresPayload, _ := json.Marshal(event.AplicaresBedPayload{
			KodeKelas:          "KL1",
			KodeRuang:          id,
			NamaRuang:          "Ruang Rawat Inap",
			Kapasitas:          1,
			Tersedia:           tersedia,
			TersediaPriaWanita:  tersedia,
		})

		task := event.BPJSTask{
			TaskID:     ulid.Make().String(),
			Action:     "APLICARES_SYNC_BED",
			ResourceID: id,
			Payload:    aplicaresPayload,
			Timestamp:  now,
		}

		if err := u.publisher.Publish(pubCtx, "bpjs.exchange", "bpjs.aplicares.bed", task); err != nil {
			slog.Warn("Failed to publish bpjs aplicares bed update event", slog.String("bed_id", id), slog.String("error", err.Error()))
		}
	}

	return nil
}

func (u *BedUsecaseImpl) GetBedSummary(ctx context.Context, wardClassID string) (*domain.BedSummary, error) {
	return u.repo.GetSummary(ctx, strings.TrimSpace(wardClassID))
}
