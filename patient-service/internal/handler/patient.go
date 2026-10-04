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
	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
)

type PatientHandler struct {
	patientv1.UnimplementedPatientServiceServer
	patient   domain.PatientUsecase
	allergy   domain.AllergyUsecase
	document  domain.DocumentUsecase
	family    domain.FamilyUsecase
	emergency domain.EmergencyContactUsecase
	insurance domain.InsuranceUsecase
}

func NewPatientHandler(
	patient domain.PatientUsecase,
	allergy domain.AllergyUsecase,
	document domain.DocumentUsecase,
	family domain.FamilyUsecase,
	emergency domain.EmergencyContactUsecase,
	insurance domain.InsuranceUsecase,
) *PatientHandler {
	return &PatientHandler{
		patient:   patient,
		allergy:   allergy,
		document:  document,
		family:    family,
		emergency: emergency,
		insurance: insurance,
	}
}

func (h *PatientHandler) CreatePatient(ctx context.Context, req *patientv1.CreatePatientRequest) (*patientv1.CreatePatientResponse, error) {
	p, err := h.patient.CreatePatient(ctx, fromProtoPatient(req.GetPatient()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.CreatePatientResponse{Patient: toProtoPatient(p)}, nil
}

func (h *PatientHandler) GetPatientByID(ctx context.Context, req *patientv1.GetPatientByIDRequest) (*patientv1.GetPatientByIDResponse, error) {
	p, err := h.patient.GetPatientByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientByIDResponse{Patient: toProtoPatient(p)}, nil
}

func (h *PatientHandler) GetPatientByNIK(ctx context.Context, req *patientv1.GetPatientByNIKRequest) (*patientv1.GetPatientByNIKResponse, error) {
	p, err := h.patient.GetPatientByNIK(ctx, req.GetNik())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientByNIKResponse{Patient: toProtoPatient(p)}, nil
}

func (h *PatientHandler) GetPatientByMedicalRecordNo(ctx context.Context, req *patientv1.GetPatientByMedicalRecordNoRequest) (*patientv1.GetPatientByMedicalRecordNoResponse, error) {
	p, err := h.patient.GetPatientByMedicalRecordNo(ctx, req.GetMedicalRecordNo())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientByMedicalRecordNoResponse{Patient: toProtoPatient(p)}, nil
}

func (h *PatientHandler) ListPatients(ctx context.Context, req *patientv1.ListPatientsRequest) (*patientv1.ListPatientsResponse, error) {
	patients, total, err := h.patient.ListPatients(ctx, domain.PatientFilter{
		Name:   req.GetName(),
		Status: domain.PatientStatus(req.GetStatus()),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.ListPatientsResponse{
		Patients:   toProtoList(patients, toProtoPatient),
		TotalCount: total,
	}, nil
}

func (h *PatientHandler) UpdatePatient(ctx context.Context, req *patientv1.UpdatePatientRequest) (*patientv1.UpdatePatientResponse, error) {
	p, err := h.patient.UpdatePatient(ctx, fromProtoPatient(req.GetPatient()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.UpdatePatientResponse{Patient: toProtoPatient(p)}, nil
}

func (h *PatientHandler) UpdatePatientStatus(ctx context.Context, req *patientv1.UpdatePatientStatusRequest) (*patientv1.UpdatePatientStatusResponse, error) {
	err := h.patient.UpdatePatientStatus(ctx, req.GetId(), domain.PatientStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.UpdatePatientStatusResponse{Success: true}, nil
}

func (h *PatientHandler) VerifyPatientBiometric(ctx context.Context, req *patientv1.VerifyPatientBiometricRequest) (*patientv1.VerifyPatientBiometricResponse, error) {
	matched, err := h.patient.VerifyPatientBiometric(ctx, req.GetPatientId(), req.GetFingerprintData())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.VerifyPatientBiometricResponse{IsMatched: matched}, nil
}

// fromProtoPatient tidak membawa status, created_at, updated_at karena diisi server.
// Enum domain memakai angka yang sama dengan enum proto, jadi cukup di-cast.
func fromProtoPatient(p *patientv1.Patient) domain.Patient {
	return domain.Patient{
		ID:                  p.GetId(),
		MedicalRecordNo:     p.GetMedicalRecordNo(),
		SatusehatID:         p.GetSatusehatId(),
		NIK:                 p.GetNik(),
		Name:                p.GetName(),
		MotherName:          p.GetMotherName(),
		BirthPlace:          p.GetBirthPlace(),
		BirthDate:           fromTimestamp(p.GetBirthDate()),
		Gender:              domain.Gender(p.GetGender()),
		BloodType:           domain.BloodType(p.GetBloodType()),
		MaritalStatus:       domain.MaritalStatus(p.GetMaritalStatus()),
		Religion:            domain.Religion(p.GetReligion()),
		Phone:               p.GetPhone(),
		Email:               p.GetEmail(),
		Address:             p.GetAddress(),
		VillageCode:         p.GetVillageCode(),
		DistrictCode:        p.GetDistrictCode(),
		CityCode:            p.GetCityCode(),
		ProvinceCode:        p.GetProvinceCode(),
		PostalCode:          p.GetPostalCode(),
		RT:                  p.GetRt(),
		RW:                  p.GetRw(),
		Occupation:          p.GetOccupation(),
		Education:           p.GetEducation(),
		PreferredLanguage:   p.GetPreferredLanguage(),
		DisabilityType:      domain.DisabilityType(p.GetDisabilityType()),
		SpecialNeedsNote:    p.GetSpecialNeedsNote(),
		PhotoURL:            p.GetPhotoUrl(),
		FingerprintTemplate: p.GetFingerprintTemplate(),
	}
}

// toProtoPatient sengaja tidak mengirim FingerprintTemplate keluar dari service.
func toProtoPatient(p *domain.Patient) *patientv1.Patient {
	return &patientv1.Patient{
		Id:                p.ID,
		MedicalRecordNo:   p.MedicalRecordNo,
		SatusehatId:       p.SatusehatID,
		Nik:               p.NIK,
		Name:              p.Name,
		MotherName:        p.MotherName,
		BirthPlace:        p.BirthPlace,
		BirthDate:         timestamppb.New(p.BirthDate),
		Gender:            patientv1.Gender(p.Gender),
		BloodType:         patientv1.BloodType(p.BloodType),
		MaritalStatus:     patientv1.MaritalStatus(p.MaritalStatus),
		Religion:          patientv1.Religion(p.Religion),
		Phone:             p.Phone,
		Email:             p.Email,
		Address:           p.Address,
		VillageCode:       p.VillageCode,
		DistrictCode:      p.DistrictCode,
		CityCode:          p.CityCode,
		ProvinceCode:      p.ProvinceCode,
		PostalCode:        p.PostalCode,
		Rt:                p.RT,
		Rw:                p.RW,
		Occupation:        p.Occupation,
		Education:         p.Education,
		PreferredLanguage: p.PreferredLanguage,
		DisabilityType:    patientv1.DisabilityType(p.DisabilityType),
		SpecialNeedsNote:  p.SpecialNeedsNote,
		PhotoUrl:          p.PhotoURL,
		Status:            patientv1.PatientStatus(p.Status),
		CreatedAt:         timestamppb.New(p.CreatedAt),
		UpdatedAt:         timestamppb.New(p.UpdatedAt),
	}
}

// toProtoList dipakai semua handler Get...(list) agar tidak mengulang loop yang sama.
func toProtoList[D, P any](items []D, convert func(*D) *P) []*P {
	out := make([]*P, len(items))
	for i := range items {
		out[i] = convert(&items[i])
	}
	return out
}

// fromTimestamp keeps a missing timestamp as the zero time,
// so the usecase rejects it instead of seeing 1970-01-01.
func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
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