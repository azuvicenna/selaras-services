package handler

import (
	"fmt"
	"net/http"
	"strconv"

	documentv1 "github.com/azuvicenna/selaras-services/document-service/gen/document/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type DocumentHandler struct {
	cm *client.ClientManager
}

func NewDocumentHandler(cm *client.ClientManager) *DocumentHandler {
	return &DocumentHandler{cm: cm}
}

// --- Templates ---

func (h *DocumentHandler) ListDocumentTemplates(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")
	typeInt, _ := strconv.Atoi(r.URL.Query().Get("type"))

	resp, err := h.cm.Document.ListDocumentTemplates(r.Context(), &documentv1.ListDocumentTemplatesRequest{
		Query:  query,
		Type:   documentv1.DocumentType(typeInt),
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Templates, page, pageSize, resp.Total)
}

// --- Documents ---

func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var req documentv1.CreateDocumentRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Document.CreateDocument(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.Document)
}

func (h *DocumentHandler) GetDocumentByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Document.GetDocumentByID(r.Context(), &documentv1.GetDocumentByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Document)
}

func (h *DocumentHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	patientID := r.URL.Query().Get("patient_id")
	encID := r.URL.Query().Get("encounter_id")
	docTypeInt, _ := strconv.Atoi(r.URL.Query().Get("type"))

	resp, err := h.cm.Document.ListDocuments(r.Context(), &documentv1.ListDocumentsRequest{
		PatientId:   patientID,
		EncounterId: encID,
		Type:        documentv1.DocumentType(docTypeInt),
		Limit:       int32(pageSize),
		Offset:      int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Documents, page, pageSize, resp.Total)
}

// Render & Download PDF directly
func (h *DocumentHandler) DownloadDocumentPDF(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Document.GenerateDocumentPDF(r.Context(), &documentv1.GenerateDocumentPDFRequest{
		DocumentId: id,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	filename := resp.Filename
	if filename == "" {
		filename = fmt.Sprintf("document-%s.pdf", id)
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(resp.PdfContent)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp.PdfContent)
}

// Verification of document QR code / signature
func (h *DocumentHandler) VerifyDocument(w http.ResponseWriter, r *http.Request) {
	docNumber := r.URL.Query().Get("number")
	hash := r.URL.Query().Get("hash")

	resp, err := h.cm.Document.VerifyDocument(r.Context(), &documentv1.VerifyDocumentRequest{
		DocumentNumber: docNumber,
		Token:          hash,
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp)
}
