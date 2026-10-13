package domain

import (
	"context"
	"time"
)

type RadiologyPriority int32

const (
	RadiologyPriorityUnspecified RadiologyPriority = 0
	RadiologyPriorityRoutine     RadiologyPriority = 1
	RadiologyPriorityUrgent      RadiologyPriority = 2
	RadiologyPriorityStat        RadiologyPriority = 3 // Cito
)

type OrderStatus int32

const (
	OrderStatusUnspecified OrderStatus = 0
	OrderStatusOrdered     OrderStatus = 1
	OrderStatusScheduled   OrderStatus = 2
	OrderStatusExamined    OrderStatus = 3
	OrderStatusCompleted   OrderStatus = 4
	OrderStatusCancelled   OrderStatus = 5
)

type RadiologyOrderItem struct {
	ID        string   `json:"id"`
	OrderID   string   `json:"order_id"`
	ExamID    string   `json:"exam_id"`
	ExamCode  string   `json:"exam_code"`
	ExamName  string   `json:"exam_name"`
	Modality  Modality `json:"modality"`
	Price     float64  `json:"price"`
	Notes     string   `json:"notes"`
}

type RadiologyOrder struct {
	ID                 string               `json:"id"`
	OrderNumber        string               `json:"order_number"`
	EncounterID        string               `json:"encounter_id"`
	PatientID          string               `json:"patient_id"`
	PatientName        string               `json:"patient_name"`
	PatientGender      string               `json:"patient_gender"`
	PatientDOB         string               `json:"patient_dob"`
	DoctorID           string               `json:"doctor_id"`
	DoctorName         string               `json:"doctor_name"`
	ClinicalIndication string               `json:"clinical_indication"`
	Priority           RadiologyPriority    `json:"priority"`
	Status             OrderStatus          `json:"status"`
	TotalPrice         float64              `json:"total_price"`
	Items              []RadiologyOrderItem `json:"items"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

type RadiologyOrderFilter struct {
	PatientID   string
	EncounterID string
	Status      OrderStatus
	Priority    RadiologyPriority
	Limit       int32
	Offset      int32
}

type RadiologyOrderRepository interface {
	Create(ctx context.Context, o *RadiologyOrder) error
	GetByID(ctx context.Context, id string) (*RadiologyOrder, error)
	List(ctx context.Context, filter RadiologyOrderFilter) ([]RadiologyOrder, int64, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) (*RadiologyOrder, error)
}
