package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
)

var (
	ErrInvalidTask         = errors.New("invalid sync task")
	ErrUnsupportedResource = errors.New("unsupported satusehat resource type")
	ErrMissingResourceID   = errors.New("missing local resource id")
	ErrSyncFailed          = errors.New("satusehat synchronization failed")
	ErrEMRSyncFailed       = errors.New("emr service sync callback failed")
)

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
)

// SyncTask defines the message received from RabbitMQ or HTTP sync trigger.
type SyncTask struct {
	TaskID       string          `json:"task_id"`
	ResourceType string          `json:"resource_type"`
	ResourceID   string          `json:"resource_id"`
	Action       ActionType      `json:"action,omitempty"`
	SatusehatID  string          `json:"satusehat_id,omitempty"`
	Payload      json.RawMessage `json:"payload,omitempty"`
	Timestamp    time.Time       `json:"timestamp"`
}

// SyncResult defines the outcome of a sync operation.
type SyncResult struct {
	TaskID       string    `json:"task_id"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	SatusehatID  string    `json:"satusehat_id"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// ToProtoResourceType maps string resource types to EMR gRPC SatusehatResourceType enum.
func ToProtoResourceType(resType string) (emrv1.SatusehatResourceType, error) {
	norm := strings.ToLower(strings.TrimSpace(resType))
	norm = strings.ReplaceAll(norm, "_", "")
	norm = strings.ReplaceAll(norm, "-", "")

	switch norm {
	case "encounter":
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ENCOUNTER, nil
	case "observation":
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_OBSERVATION, nil
	case "clinicalnote", "composition":
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_CLINICAL_NOTE, nil
	case "condition", "diagnosis":
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_CONDITION, nil
	case "procedure":
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_PROCEDURE, nil
	case "medicationrequest", "medication", "prescription":
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_MEDICATION_REQUEST, nil
	case "diagnosticreport", "labreport", "radiologyreport":
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_DIAGNOSTIC_REPORT, nil
	case "allergy", "allergyintolerance":
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ALLERGY, nil
	default:
		return emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_UNSPECIFIED, ErrUnsupportedResource
	}
}

// ToFHIRResourceType returns the official HL7 FHIR R4 resource name used in URL paths.
func ToFHIRResourceType(resType string) (string, error) {
	protoType, err := ToProtoResourceType(resType)
	if err != nil {
		return "", err
	}

	switch protoType {
	case emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ENCOUNTER:
		return "Encounter", nil
	case emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_OBSERVATION:
		return "Observation", nil
	case emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_CLINICAL_NOTE:
		return "Composition", nil
	case emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_CONDITION:
		return "Condition", nil
	case emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_PROCEDURE:
		return "Procedure", nil
	case emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_MEDICATION_REQUEST:
		return "MedicationRequest", nil
	case emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_DIAGNOSTIC_REPORT:
		return "DiagnosticReport", nil
	case emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ALLERGY:
		return "AllergyIntolerance", nil
	default:
		return "", ErrUnsupportedResource
	}
}
