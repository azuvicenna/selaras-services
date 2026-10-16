package domain_test

import (
	"testing"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/domain"
)

func TestToProtoResourceType(t *testing.T) {
	tests := []struct {
		input       string
		expected    emrv1.SatusehatResourceType
		shouldError bool
	}{
		{"encounter", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ENCOUNTER, false},
		{"Encounter", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ENCOUNTER, false},
		{"observation", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_OBSERVATION, false},
		{"clinical_note", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_CLINICAL_NOTE, false},
		{"composition", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_CLINICAL_NOTE, false},
		{"condition", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_CONDITION, false},
		{"diagnosis", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_CONDITION, false},
		{"procedure", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_PROCEDURE, false},
		{"medication_request", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_MEDICATION_REQUEST, false},
		{"prescription", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_MEDICATION_REQUEST, false},
		{"diagnostic_report", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_DIAGNOSTIC_REPORT, false},
		{"allergy", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ALLERGY, false},
		{"allergy_intolerance", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_ALLERGY, false},
		{"unknown_resource", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_UNSPECIFIED, true},
		{"", emrv1.SatusehatResourceType_SATUSEHAT_RESOURCE_TYPE_UNSPECIFIED, true},
	}

	for _, tt := range tests {
		got, err := domain.ToProtoResourceType(tt.input)
		if tt.shouldError {
			if err == nil {
				t.Errorf("expected error for input %q, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("expected %v, got %v for input %q", tt.expected, got, tt.input)
			}
		}
	}
}

func TestToFHIRResourceType(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		shouldError bool
	}{
		{"encounter", "Encounter", false},
		{"observation", "Observation", false},
		{"clinical_note", "Composition", false},
		{"condition", "Condition", false},
		{"procedure", "Procedure", false},
		{"medication_request", "MedicationRequest", false},
		{"diagnostic_report", "DiagnosticReport", false},
		{"allergy", "AllergyIntolerance", false},
		{"unknown", "", true},
	}

	for _, tt := range tests {
		got, err := domain.ToFHIRResourceType(tt.input)
		if tt.shouldError {
			if err == nil {
				t.Errorf("expected error for input %q, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("expected %v, got %v for input %q", tt.expected, got, tt.input)
			}
		}
	}
}
