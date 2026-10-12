package domain

import (
	"context"
	"time"
)

type BedStatus int32

const (
	BedStatusUnspecified BedStatus = 0
	BedStatusAvailable   BedStatus = 1
	BedStatusOccupied    BedStatus = 2
	BedStatusReserved    BedStatus = 3
	BedStatusCleaning    BedStatus = 4
	BedStatusMaintenance BedStatus = 5
	BedStatusIsolation   BedStatus = 6
)

type Bed struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	WardID    string    `json:"ward_id"`
	Status    BedStatus `json:"status"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BedFilter struct {
	Query    string
	WardID   string
	Status   BedStatus
	IsActive *bool
	Limit    int
	Offset   int
}

type BedSummaryItem struct {
	WardClassID     string `json:"ward_class_id"`
	WardClassName   string `json:"ward_class_name"`
	WardClassCode   string `json:"ward_class_code"`
	TotalBeds       int32  `json:"total_beds"`
	AvailableBeds   int32  `json:"available_beds"`
	OccupiedBeds    int32  `json:"occupied_beds"`
	ReservedBeds    int32  `json:"reserved_beds"`
	CleaningBeds    int32  `json:"cleaning_beds"`
	MaintenanceBeds int32  `json:"maintenance_beds"`
	IsolationBeds   int32  `json:"isolation_beds"`
}

type BedSummary struct {
	Items           []BedSummaryItem `json:"items"`
	TotalBeds       int32            `json:"total_beds"`
	AvailableBeds   int32            `json:"available_beds"`
	OccupiedBeds    int32            `json:"occupied_beds"`
	ReservedBeds    int32            `json:"reserved_beds"`
	CleaningBeds    int32            `json:"cleaning_beds"`
	MaintenanceBeds int32            `json:"maintenance_beds"`
	IsolationBeds   int32            `json:"isolation_beds"`
}

type BedRepository interface {
	Create(ctx context.Context, item *Bed) error
	GetByID(ctx context.Context, id string) (*Bed, error)
	GetByCode(ctx context.Context, wardID, code string) (*Bed, error)
	List(ctx context.Context, filter BedFilter) ([]Bed, int64, error)
	Update(ctx context.Context, item *Bed) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status BedStatus, updatedAt time.Time) error
	GetSummary(ctx context.Context, wardClassID string) (*BedSummary, error)
}

type BedUsecase interface {
	CreateBed(ctx context.Context, item Bed) (*Bed, error)
	GetBedByID(ctx context.Context, id string) (*Bed, error)
	GetBedByCode(ctx context.Context, wardID, code string) (*Bed, error)
	ListBeds(ctx context.Context, filter BedFilter) ([]Bed, int64, error)
	UpdateBed(ctx context.Context, item Bed) (*Bed, error)
	DeleteBed(ctx context.Context, id string) error
	UpdateBedStatus(ctx context.Context, id string, status BedStatus) error
	GetBedSummary(ctx context.Context, wardClassID string) (*BedSummary, error)
}
