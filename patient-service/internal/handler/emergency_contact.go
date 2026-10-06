package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

func (h *PatientHandler) AddEmergencyContact(ctx context.Context, req *patientv1.AddEmergencyContactRequest) (*patientv1.AddEmergencyContactResponse, error) {
	c, err := h.emergency.AddEmergencyContact(ctx, fromProtoEmergencyContact(req.GetContact()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.AddEmergencyContactResponse{Contact: toProtoEmergencyContact(c)}, nil
}

func (h *PatientHandler) GetEmergencyContacts(ctx context.Context, req *patientv1.GetEmergencyContactsRequest) (*patientv1.GetEmergencyContactsResponse, error) {
	contacts, err := h.emergency.GetEmergencyContacts(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetEmergencyContactsResponse{Contacts: toProtoList(contacts, toProtoEmergencyContact)}, nil
}

func (h *PatientHandler) UpdateEmergencyContact(ctx context.Context, req *patientv1.UpdateEmergencyContactRequest) (*patientv1.UpdateEmergencyContactResponse, error) {
	c, err := h.emergency.UpdateEmergencyContact(ctx, fromProtoEmergencyContact(req.GetContact()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.UpdateEmergencyContactResponse{Contact: toProtoEmergencyContact(c)}, nil
}

func (h *PatientHandler) DeleteEmergencyContact(ctx context.Context, req *patientv1.DeleteEmergencyContactRequest) (*patientv1.DeleteEmergencyContactResponse, error) {
	if err := h.emergency.DeleteEmergencyContact(ctx, req.GetId()); err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.DeleteEmergencyContactResponse{Success: true}, nil
}

func fromProtoEmergencyContact(c *patientv1.PatientEmergencyContact) domain.PatientEmergencyContact {
	return domain.PatientEmergencyContact{
		ID:           c.GetId(),
		PatientID:    c.GetPatientId(),
		Name:         c.GetName(),
		Relationship: c.GetRelationship(),
		Phone:        c.GetPhone(),
		Address:      c.GetAddress(),
	}
}

func toProtoEmergencyContact(c *domain.PatientEmergencyContact) *patientv1.PatientEmergencyContact {
	return &patientv1.PatientEmergencyContact{
		Id:           c.ID,
		PatientId:    c.PatientID,
		Name:         c.Name,
		Relationship: c.Relationship,
		Phone:        c.Phone,
		Address:      c.Address,
		CreatedAt:    timestamppb.New(c.CreatedAt),
		UpdatedAt:    timestamppb.New(c.UpdatedAt),
	}
}
