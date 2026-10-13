package domain

import (
	"context"
	"time"
)

type LabPriority int32

const (
	LabPriorityUnspecified LabPriority = 0
	LabPriorityRoutine     LabPriority = 1
	LabPriorityUrgent      LabPriority = 2
	LabPriorityStat        LabPriority = 3 // Cito
)

type OrderStatus int32

const (
	OrderStatusUnspecified       OrderStatus = 0
	OrderStatusOrdered           OrderStatus = 1
	OrderStatusSampleCollected   OrderStatus = 2
	OrderStatusInProgress        OrderStatus = 3
	OrderStatusCompleted         OrderStatus = 4
	OrderStatusCancelled         OrderStatus = 5
)

type LabOrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	TestID    string  `json:"test_id"`
	TestCode  string  `json:"test_code"`
	TestName  string  `json:"test_name"`
	Price     float64 `json:"price"`
	Notes     string  `json:"notes"`
}

type LabOrder struct {
	ID            string         `json:"id"`
	OrderNumber   string         `json:"order_number"`
	EncounterID   string         `json:"encounter_id"`
	PatientID     string         `json:"patient_id"`
	PatientName   string         `json:"patient_name"`
	PatientGender string         `json:"patient_gender"`
	PatientDOB    string         `json:"patient_dob"`
	DoctorID      string         `json:"doctor_id"`
	DoctorName    string         `json:"doctor_name"`
	ClinicalNotes string         `json:"clinical_notes"`
	Priority      LabPriority    `json:"priority"`
	Status        OrderStatus    `json:"status"`
	TotalPrice    float64        `json:"total_price"`
	Items         []LabOrderItem `json:"items"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type LabOrderFilter struct {
	PatientID   string
	EncounterID string
	Status      OrderStatus
	Priority    LabPriority
	Limit       int
	Offset      int
}

type LabOrderRepository interface {
	Create(ctx context.Context, o *LabOrder) error
	GetByID(ctx context.Context, id string) (*LabOrder, error)
	List(ctx context.Context, filter LabOrderFilter) ([]LabOrder, int64, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
}
