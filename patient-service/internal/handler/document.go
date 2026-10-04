package handler

import (
	"bytes"
	"context"
	"io"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
)

func (h *PatientHandler) UploadPatientDocument(ctx context.Context, req *patientv1.UploadPatientDocumentRequest) (*patientv1.UploadPatientDocumentResponse, error) {
	// Reader dibiarkan nil kalau file kosong, supaya validasi usecase yang menolaknya.
	var file io.Reader
	if len(req.GetFileContent()) > 0 {
		file = bytes.NewReader(req.GetFileContent())
	}

	d, err := h.document.UploadPatientDocument(
		ctx,
		fromProtoDocument(req.GetDocument()),
		file,
		req.GetFileName(),
		req.GetContentType(),
	)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.UploadPatientDocumentResponse{Document: toProtoDocument(d)}, nil
}

func (h *PatientHandler) GetPatientDocuments(ctx context.Context, req *patientv1.GetPatientDocumentsRequest) (*patientv1.GetPatientDocumentsResponse, error) {
	docs, err := h.document.GetPatientDocuments(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientDocumentsResponse{Documents: toProtoList(docs, toProtoDocument)}, nil
}

func (h *PatientHandler) GetPatientDocumentByID(ctx context.Context, req *patientv1.GetPatientDocumentByIDRequest) (*patientv1.GetPatientDocumentByIDResponse, error) {
	d, err := h.document.GetPatientDocumentByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.GetPatientDocumentByIDResponse{Document: toProtoDocument(d)}, nil
}

func (h *PatientHandler) DeletePatientDocument(ctx context.Context, req *patientv1.DeletePatientDocumentRequest) (*patientv1.DeletePatientDocumentResponse, error) {
	if err := h.document.DeletePatientDocument(ctx, req.GetId()); err != nil {
		return nil, toStatusError(err)
	}
	return &patientv1.DeletePatientDocumentResponse{Success: true}, nil
}

// file_url diisi server setelah upload, jadi tidak dibaca dari request.
func fromProtoDocument(d *patientv1.PatientDocument) domain.PatientDocument {
	return domain.PatientDocument{
		ID:             d.GetId(),
		PatientID:      d.GetPatientId(),
		Type:           domain.DocumentType(d.GetType()),
		DocumentNumber: d.GetDocumentNumber(),
	}
}

func toProtoDocument(d *domain.PatientDocument) *patientv1.PatientDocument {
	return &patientv1.PatientDocument{
		Id:             d.ID,
		PatientId:      d.PatientID,
		Type:           patientv1.DocumentType(d.Type),
		DocumentNumber: d.DocumentNumber,
		FileUrl:        d.FileURL,
		CreatedAt:      timestamppb.New(d.CreatedAt),
		UpdatedAt:      timestamppb.New(d.UpdatedAt),
	}
}