package domain

import (
	"context"
	"time"
)

type RadiologyReport struct {
	ID              string     `json:"id"`
	OrderID         string     `json:"order_id"`
	OrderNumber     string     `json:"order_number"`
	PatientID       string     `json:"patient_id"`
	RadiologistID   string     `json:"radiologist_id"`
	RadiologistName string     `json:"radiologist_name"`
	ClinicalHistory string     `json:"clinical_history"`
	Findings        string     `json:"findings"`
	Impression      string     `json:"impression"`
	Recommendation  string     `json:"recommendation"`
	IsCritical      bool       `json:"is_critical"`
	IsVerified      bool       `json:"is_verified"`
	VerifiedAt      *time.Time `json:"verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type RadiologyReportRepository interface {
	Save(ctx context.Context, r *RadiologyReport) error
	GetByOrderID(ctx context.Context, orderID string) (*RadiologyReport, error)
	Verify(ctx context.Context, orderID, radiologistID, radiologistName string) (*RadiologyReport, error)
}
