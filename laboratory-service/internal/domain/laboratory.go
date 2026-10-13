package domain

import (
	"context"
)

type LaboratoryUsecase interface {
	// Master Pemeriksaan Lab
	CreateLabTest(ctx context.Context, t LabTest) (*LabTest, error)
	GetLabTestByID(ctx context.Context, id string) (*LabTest, error)
	GetLabTestByCode(ctx context.Context, code string) (*LabTest, error)
	ListLabTests(ctx context.Context, filter LabTestFilter) ([]LabTest, int64, error)
	UpdateLabTest(ctx context.Context, t LabTest) (*LabTest, error)
	DeleteLabTest(ctx context.Context, id string) error

	// Order Lab
	CreateLabOrder(ctx context.Context, o LabOrder, testIDs []string, notes []string) (*LabOrder, error)
	GetLabOrderByID(ctx context.Context, id string) (*LabOrder, error)
	ListLabOrders(ctx context.Context, filter LabOrderFilter) ([]LabOrder, int64, error)
	UpdateLabOrderStatus(ctx context.Context, id string, status OrderStatus) (*LabOrder, error)

	// Spesimen Lab
	CollectSpecimen(ctx context.Context, s LabSpecimen) (*LabSpecimen, error)
	GetSpecimensByOrderID(ctx context.Context, orderID string) ([]LabSpecimen, error)
	UpdateSpecimenStatus(ctx context.Context, id string, status SpecimenStatus, notes string) (*LabSpecimen, error)

	// Hasil Lab & Verifikasi
	SaveLabResult(ctx context.Context, orderID, analystID, analystName, overallNotes string, items []LabResultItem) (*LabResult, error)
	GetLabResultByOrderID(ctx context.Context, orderID string) (*LabResult, error)
	VerifyLabResult(ctx context.Context, orderID, pathologistID, pathologistName, notes string) (*LabResult, error)
}
