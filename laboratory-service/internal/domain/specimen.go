package domain

import (
	"context"
	"time"
)

type SpecimenStatus int32

const (
	SpecimenStatusUnspecified SpecimenStatus = 0
	SpecimenStatusCollected   SpecimenStatus = 1
	SpecimenStatusReceived    SpecimenStatus = 2
	SpecimenStatusRejected    SpecimenStatus = 3
)

type LabSpecimen struct {
	ID             string         `json:"id"`
	OrderID        string         `json:"order_id"`
	SpecimenNumber string         `json:"specimen_number"`
	SpecimenType   SpecimenType   `json:"specimen_type"`
	Status         SpecimenStatus `json:"status"`
	CollectorID    string         `json:"collector_id"`
	CollectorName  string         `json:"collector_name"`
	Notes          string         `json:"notes"`
	CollectedAt    *time.Time     `json:"collected_at,omitempty"`
	ReceivedAt     *time.Time     `json:"received_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type LabSpecimenRepository interface {
	Create(ctx context.Context, s *LabSpecimen) error
	GetByID(ctx context.Context, id string) (*LabSpecimen, error)
	GetByOrderID(ctx context.Context, orderID string) ([]LabSpecimen, error)
	UpdateStatus(ctx context.Context, id string, status SpecimenStatus, notes string, receivedAt *time.Time) (*LabSpecimen, error)
}
