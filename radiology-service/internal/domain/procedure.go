package domain

import (
	"context"
	"time"
)

type RadiologyProcedure struct {
	ID                string     `json:"id"`
	OrderID           string     `json:"order_id"`
	RoomOrDevice      string     `json:"room_or_device"`
	RadiographerID    string     `json:"radiographer_id"`
	RadiographerName  string     `json:"radiographer_name"`
	ContrastAgentUsed string     `json:"contrast_agent_used"`
	SeriesCount       int32      `json:"series_count"`
	InstanceCount     int32      `json:"instance_count"`
	DICOMStudyUID     string     `json:"dicom_study_uid"`
	PACSViewerURL     string     `json:"pacs_viewer_url"`
	Notes             string     `json:"notes"`
	ScheduledAt       *time.Time `json:"scheduled_at"`
	PerformedAt       *time.Time `json:"performed_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type RadiologyProcedureRepository interface {
	Schedule(ctx context.Context, p *RadiologyProcedure) error
	Record(ctx context.Context, p *RadiologyProcedure) error
	GetByOrderID(ctx context.Context, orderID string) (*RadiologyProcedure, error)
}
