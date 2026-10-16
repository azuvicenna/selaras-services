package syncer_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/syncer"
)

type mockSatusehatClient struct {
	postResultID string
	postError    error
	putResultID  string
	putError     error
	lastPayload  any
}

func (m *mockSatusehatClient) PostResource(ctx context.Context, fhirType string, payload any) (string, error) {
	m.lastPayload = payload
	return m.postResultID, m.postError
}

func (m *mockSatusehatClient) PutResource(ctx context.Context, fhirType, id string, payload any) (string, error) {
	m.lastPayload = payload
	return m.putResultID, m.putError
}

func (m *mockSatusehatClient) GetResource(ctx context.Context, fhirType, id string) ([]byte, error) {
	return nil, nil
}

type mockEMRClient struct {
	lastResourceType emrv1.SatusehatResourceType
	lastLocalID      string
	lastSatusehatID  string
	syncError        error
}

func (m *mockEMRClient) SyncSatusehatID(ctx context.Context, resType emrv1.SatusehatResourceType, localID, satusehatID string) error {
	m.lastResourceType = resType
	m.lastLocalID = localID
	m.lastSatusehatID = satusehatID
	return m.syncError
}

func (m *mockEMRClient) Close() error {
	return nil
}

func TestSyncer_Process_SuccessCreate(t *testing.T) {
	mockSS := &mockSatusehatClient{
		postResultID: "ENC-IHS-98765",
	}
	mockEMR := &mockEMRClient{}

	s := syncer.NewSyncer(mockSS, mockEMR, "10000004", nil)

	task := domain.SyncTask{
		TaskID:       "task-01",
		ResourceType: "encounter",
		ResourceID:   "local-enc-123",
		Action:       domain.ActionCreate,
		Payload: json.RawMessage(`{
			"patient_ihs": "10000000001",
			"status": "finished",
			"start_time": "2026-10-16T08:00:00+07:00"
		}`),
	}

	result, err := s.Process(context.Background(), task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success to be true")
	}
	if result.SatusehatID != "ENC-IHS-98765" {
		t.Errorf("expected SatusehatID ENC-IHS-98765, got %s", result.SatusehatID)
	}

	// Verify EMR gRPC callback arguments
	if mockEMR.lastResourceType != emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ENCOUNTER {
		t.Errorf("expected resource type ENCOUNTER, got %v", mockEMR.lastResourceType)
	}
	if mockEMR.lastLocalID != "local-enc-123" {
		t.Errorf("expected local ID local-enc-123, got %s", mockEMR.lastLocalID)
	}
	if mockEMR.lastSatusehatID != "ENC-IHS-98765" {
		t.Errorf("expected satusehat ID ENC-IHS-98765, got %s", mockEMR.lastSatusehatID)
	}
}

func TestSyncer_Process_SatusehatError(t *testing.T) {
	mockSS := &mockSatusehatClient{
		postError: errors.New("network connection timeout"),
	}
	mockEMR := &mockEMRClient{}

	s := syncer.NewSyncer(mockSS, mockEMR, "10000004", nil)

	task := domain.SyncTask{
		TaskID:       "task-02",
		ResourceType: "condition",
		ResourceID:   "local-cond-1",
		Payload: json.RawMessage(`{
			"patient_ihs": "10000000001",
			"encounter_ihs": "ENC-001",
			"icd10_code": "A09",
			"diagnosis": "Gastroenteritis"
		}`),
	}

	result, err := s.Process(context.Background(), task)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result.Success {
		t.Error("expected result.Success to be false")
	}
}

func TestSyncer_Process_EMRCallbackError(t *testing.T) {
	mockSS := &mockSatusehatClient{
		postResultID: "OBS-IHS-12345",
	}
	mockEMR := &mockEMRClient{
		syncError: errors.New("emr database update failed"),
	}

	s := syncer.NewSyncer(mockSS, mockEMR, "10000004", nil)

	task := domain.SyncTask{
		TaskID:       "task-03",
		ResourceType: "observation",
		ResourceID:   "local-obs-1",
		Payload: json.RawMessage(`{
			"patient_ihs": "10000000001",
			"encounter_ihs": "ENC-001",
			"loinc_code": "8867-4",
			"name": "Heart rate"
		}`),
	}

	result, err := s.Process(context.Background(), task)
	if err == nil {
		t.Fatal("expected error from EMR callback failure, got nil")
	}
	if result.Success {
		t.Error("expected result.Success to be false")
	}
}

func TestSyncer_Process_MissingResourceID(t *testing.T) {
	s := syncer.NewSyncer(&mockSatusehatClient{}, &mockEMRClient{}, "10000004", nil)

	task := domain.SyncTask{
		TaskID:       "task-04",
		ResourceType: "encounter",
		ResourceID:   "", // Missing
	}

	_, err := s.Process(context.Background(), task)
	if !errors.Is(err, domain.ErrMissingResourceID) {
		t.Errorf("expected ErrMissingResourceID, got %v", err)
	}
}
