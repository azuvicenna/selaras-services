package domain

import (
	"context"
)

type RadiologyUsecase interface {
	// Master Pemeriksaan Radiologi
	CreateRadiologyExam(ctx context.Context, e RadiologyExam) (*RadiologyExam, error)
	GetRadiologyExamByID(ctx context.Context, id string) (*RadiologyExam, error)
	GetRadiologyExamByCode(ctx context.Context, code string) (*RadiologyExam, error)
	ListRadiologyExams(ctx context.Context, filter RadiologyExamFilter) ([]RadiologyExam, int64, error)
	UpdateRadiologyExam(ctx context.Context, e RadiologyExam) (*RadiologyExam, error)
	DeleteRadiologyExam(ctx context.Context, id string) error

	// Order Radiologi
	CreateRadiologyOrder(ctx context.Context, o RadiologyOrder, examIDs []string, notes []string) (*RadiologyOrder, error)
	GetRadiologyOrderByID(ctx context.Context, id string) (*RadiologyOrder, error)
	ListRadiologyOrders(ctx context.Context, filter RadiologyOrderFilter) ([]RadiologyOrder, int64, error)
	UpdateRadiologyOrderStatus(ctx context.Context, id string, status OrderStatus) (*RadiologyOrder, error)

	// Penjadwalan & Pelaksanaan Prosedur Radiologi
	ScheduleRadiologyProcedure(ctx context.Context, p RadiologyProcedure) (*RadiologyProcedure, error)
	RecordRadiologyProcedure(ctx context.Context, p RadiologyProcedure) (*RadiologyProcedure, error)
	GetRadiologyProcedureByOrderID(ctx context.Context, orderID string) (*RadiologyProcedure, error)

	// Ekspertise & Pelaporan Radiologi
	SaveRadiologyReport(ctx context.Context, r RadiologyReport) (*RadiologyReport, error)
	GetRadiologyReportByOrderID(ctx context.Context, orderID string) (*RadiologyReport, error)
	VerifyRadiologyReport(ctx context.Context, orderID, radiologistID, radiologistName string) (*RadiologyReport, error)
}
