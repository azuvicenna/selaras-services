package domain

import (
	"context"
	"time"
)

type ResultFlag int32

const (
	ResultFlagUnspecified  ResultFlag = 0
	ResultFlagNormal       ResultFlag = 1
	ResultFlagLow          ResultFlag = 2
	ResultFlagHigh         ResultFlag = 3
	ResultFlagCriticalLow  ResultFlag = 4
	ResultFlagCriticalHigh ResultFlag = 5
	ResultFlagAbnormal     ResultFlag = 6
)

type LabResultItem struct {
	ID             string     `json:"id"`
	ResultID       string     `json:"result_id"`
	OrderID        string     `json:"order_id"`
	TestID         string     `json:"test_id"`
	TestName       string     `json:"test_name"`
	LOINCCode      string     `json:"loinc_code"`
	ResultValue    string     `json:"result_value"`
	NumericValue   float64    `json:"numeric_value"`
	ReferenceRange string     `json:"reference_range"`
	Unit           string     `json:"unit"`
	Flag           ResultFlag `json:"flag"`
	Notes          string     `json:"notes"`
}

type LabResult struct {
	ID              string          `json:"id"`
	OrderID         string          `json:"order_id"`
	OrderNumber     string          `json:"order_number"`
	PatientID       string          `json:"patient_id"`
	AnalystID       string          `json:"analyst_id"`
	AnalystName     string          `json:"analyst_name"`
	PathologistID   string          `json:"pathologist_id"`
	PathologistName string          `json:"pathologist_name"`
	IsVerified      bool            `json:"is_verified"`
	OverallNotes    string          `json:"overall_notes"`
	Items           []LabResultItem `json:"items"`
	VerifiedAt      *time.Time      `json:"verified_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type LabResultRepository interface {
	SaveResult(ctx context.Context, r *LabResult) error
	GetByOrderID(ctx context.Context, orderID string) (*LabResult, error)
	VerifyResult(ctx context.Context, orderID string, pathologistID, pathologistName, notes string, verifiedAt time.Time) (*LabResult, error)
}
