package domain

import (
	"context"
	"time"
)

type AdmissionType int32

const (
	AdmissionTypeUnspecified AdmissionType = 0
	AdmissionTypeOutpatient  AdmissionType = 1
	AdmissionTypeInpatient   AdmissionType = 2
	AdmissionTypeEmergency   AdmissionType = 3
)

type AdmissionStatus int32

const (
	AdmissionStatusUnspecified AdmissionStatus = 0
	AdmissionStatusRegistered  AdmissionStatus = 1
	AdmissionStatusAdmitted    AdmissionStatus = 2
	AdmissionStatusDischarged  AdmissionStatus = 3
	AdmissionStatusCancelled   AdmissionStatus = 4
)

type Admission struct {
	ID              string          `json:"id"`
	AdmissionNumber string          `json:"admission_number"`
	PatientID       string          `json:"patient_id"`
	AppointmentID   *string         `json:"appointment_id,omitempty"`
	AdmissionType   AdmissionType   `json:"admission_type"`
	Status          AdmissionStatus `json:"status"`
	DepartmentID    string          `json:"department_id"`
	PractitionerID  string          `json:"practitioner_id"`
	AdmissionDate   time.Time       `json:"admission_date"`
	DischargeDate   *time.Time      `json:"discharge_date,omitempty"`
	WardID          *string         `json:"ward_id,omitempty"`
	BedID           *string         `json:"bed_id,omitempty"`
	Notes           string          `json:"notes"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type AdmissionFilter struct {
	PatientID     string
	DepartmentID  string
	AdmissionType AdmissionType
	Status        AdmissionStatus
	StartDate     string
	EndDate       string
	Limit         int
	Offset        int
}

type AdmissionRepository interface {
	Create(ctx context.Context, item *Admission) error
	GetByID(ctx context.Context, id string) (*Admission, error)
	GetByNumber(ctx context.Context, number string) (*Admission, error)
	List(ctx context.Context, filter AdmissionFilter) ([]Admission, int64, error)
	Update(ctx context.Context, item *Admission) error
	UpdateStatus(ctx context.Context, id string, status AdmissionStatus, dischargeDate *time.Time) (*Admission, error)
	Delete(ctx context.Context, id string) error
}

type AdmissionUsecase interface {
	CreateAdmission(ctx context.Context, item Admission) (*Admission, error)
	GetAdmissionByID(ctx context.Context, id string) (*Admission, error)
	GetAdmissionByNumber(ctx context.Context, number string) (*Admission, error)
	ListAdmissions(ctx context.Context, filter AdmissionFilter) ([]Admission, int64, error)
	UpdateAdmission(ctx context.Context, item Admission) (*Admission, error)
	UpdateAdmissionStatus(ctx context.Context, id string, status AdmissionStatus, dischargeDate *time.Time) (*Admission, error)
	DeleteAdmission(ctx context.Context, id string) error
}
