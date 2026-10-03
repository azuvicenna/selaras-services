package handler

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	"github.com/azuvicenna/selaras-services/patient-service/internal/usecase"
	patientv1 "github.com/azuvicenna/selaras-services/patient-service/proto/patient/v1"
)

type PatientHandler struct {
	patientv1.UnimplementedPatientServiceServer
	uc *usecase.PatientUsecase
}

func NewPatientHandler(uc *usecase.PatientUsecase) *PatientHandler {
	return &PatientHandler{uc: uc}
}

func (h *PatientHandler) RegisterPatient(ctx context.Context, req *patientv1.RegisterPatientRequest) (*patientv1.RegisterPatientResponse, error) {
	p, err := h.uc.Register(ctx, domain.Patient{
		NIK:       req.GetNik(),
		Name:      req.GetName(),
		BirthDate: fromTimestamp(req.GetBirthDate()),
		Gender:    toDomainGender(req.GetGender()),
		Phone:     req.GetPhone(),
		Address:   req.GetAddress(),
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.RegisterPatientResponse{Patient: toProto(p)}, nil
}

func (h *PatientHandler) GetPatient(ctx context.Context, req *patientv1.GetPatientRequest) (*patientv1.GetPatientResponse, error) {
	p, err := h.uc.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientResponse{Patient: toProto(p)}, nil
}

func (h *PatientHandler) ListPatients(ctx context.Context, req *patientv1.ListPatientsRequest) (*patientv1.ListPatientsResponse, error) {
	patients, err := h.uc.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, toStatusError(err)
	}
	res := &patientv1.ListPatientsResponse{Patients: make([]*patientv1.Patient, 0, len(patients))}
	for i := range patients {
		res.Patients = append(res.Patients, toProto(&patients[i]))
	}
	return res, nil
}

func (h *PatientHandler) UpdatePatient(ctx context.Context, req *patientv1.UpdatePatientRequest) (*patientv1.UpdatePatientResponse, error) {
	p, err := h.uc.Update(ctx, domain.Patient{
		ID:        req.GetId(),
		NIK:       req.GetNik(),
		Name:      req.GetName(),
		BirthDate: fromTimestamp(req.GetBirthDate()),
		Gender:    toDomainGender(req.GetGender()),
		Phone:     req.GetPhone(),
		Address:   req.GetAddress(),
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.UpdatePatientResponse{Patient: toProto(p)}, nil
}

func toProto(p *domain.Patient) *patientv1.Patient {
	return &patientv1.Patient{
		Id:              p.ID,
		MedicalRecordNo: p.MedicalRecordNo,
		Nik:             p.NIK,
		Name:            p.Name,
		BirthDate:       timestamppb.New(p.BirthDate),
		Gender:          toProtoGender(p.Gender),
		Phone:           p.Phone,
		Address:         p.Address,
		CreatedAt:       timestamppb.New(p.CreatedAt),
		UpdatedAt:       timestamppb.New(p.UpdatedAt),
	}
}

// fromTimestamp keeps a missing timestamp as the zero time,
// so the usecase rejects it instead of seeing 1970-01-01.
func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func toDomainGender(g patientv1.Gender) domain.Gender {
	switch g {
	case patientv1.Gender_GENDER_MALE:
		return domain.Male
	case patientv1.Gender_GENDER_FEMALE:
		return domain.Female
	}
	return ""
}

func toProtoGender(g domain.Gender) patientv1.Gender {
	switch g {
	case domain.Male:
		return patientv1.Gender_GENDER_MALE
	case domain.Female:
		return patientv1.Gender_GENDER_FEMALE
	}
	return patientv1.Gender_GENDER_UNSPECIFIED
}

func toStatusError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	}
	slog.Error("unexpected error", "err", err)
	return status.Error(codes.Internal, "internal error")
}