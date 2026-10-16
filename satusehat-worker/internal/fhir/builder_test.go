package fhir_test

import (
	"encoding/json"
	"testing"

	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/fhir"
)

func TestBuildEncounter(t *testing.T) {
	enc := fhir.BuildEncounter("10000004", fhir.EncounterParams{
		LocalID:          "ENC-001",
		PatientIHS:       "10000000001",
		PatientName:      "Budi Santoso",
		PractitionerIHS:  "N10000001",
		PractitionerName: "dr. Andi",
		LocationID:       "LOC-POLI-01",
		LocationName:     "Poli Umum",
		Status:           "finished",
		ClassCode:        "AMB",
		StartTime:        "2026-10-16T08:00:00+07:00",
		EndTime:          "2026-10-16T08:30:00+07:00",
	})

	if enc.ResourceType != "Encounter" {
		t.Fatalf("expected Encounter, got %s", enc.ResourceType)
	}
	if enc.Subject.Reference != "Patient/10000000001" {
		t.Errorf("expected Patient/10000000001, got %s", enc.Subject.Reference)
	}
	if len(enc.Identifier) == 0 || enc.Identifier[0].Value != "ENC-001" {
		t.Errorf("expected identifier ENC-001")
	}
	if len(enc.Participant) == 0 || enc.Participant[0].Individual.Reference != "Practitioner/N10000001" {
		t.Errorf("expected practitioner reference Practitioner/N10000001")
	}
}

func TestBuildCondition(t *testing.T) {
	cond := fhir.BuildCondition(fhir.ConditionParams{
		PatientIHS:   "10000000001",
		EncounterIHS: "ENC-FHIR-999",
		ICD10Code:    "A09",
		Diagnosis:    "Infectious gastroenteritis",
	})

	if cond.ResourceType != "Condition" {
		t.Fatalf("expected Condition, got %s", cond.ResourceType)
	}
	if cond.Code.Coding[0].Code != "A09" {
		t.Errorf("expected A09, got %s", cond.Code.Coding[0].Code)
	}
	if cond.Encounter.Reference != "Encounter/ENC-FHIR-999" {
		t.Errorf("expected Encounter/ENC-FHIR-999, got %s", cond.Encounter.Reference)
	}
}

func TestBuildObservation(t *testing.T) {
	val := 120.0
	obs := fhir.BuildObservation(fhir.ObservationParams{
		PatientIHS:    "10000000001",
		EncounterIHS:  "ENC-FHIR-999",
		LOINCCode:     "8480-6",
		Name:          "Systolic blood pressure",
		Category:      "vital-signs",
		ValueQuantity: &val,
		Unit:          "mm[Hg]",
	})

	if obs.ResourceType != "Observation" {
		t.Fatalf("expected Observation, got %s", obs.ResourceType)
	}
	if obs.ValueQuantity == nil || obs.ValueQuantity.Value != 120.0 {
		t.Errorf("expected value 120.0, got %v", obs.ValueQuantity)
	}
}

func TestPrepareFHIRPayload_PassthroughRawFHIR(t *testing.T) {
	raw := json.RawMessage(`{"resourceType":"Encounter","status":"finished"}`)
	res, err := fhir.PrepareFHIRPayload("10000004", "encounter", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(res.(json.RawMessage)) != string(raw) {
		t.Errorf("expected passthrough of raw FHIR, got %v", res)
	}
}

func TestPrepareFHIRPayload_BuildFromParams(t *testing.T) {
	params := json.RawMessage(`{
		"patient_ihs": "10000000001",
		"encounter_ihs": "ENC-123",
		"icd10_code": "J00",
		"diagnosis": "Nasopharyngitis"
	}`)

	res, err := fhir.PrepareFHIRPayload("10000004", "condition", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cond, ok := res.(*fhir.ConditionResource)
	if !ok {
		t.Fatalf("expected *ConditionResource, got %T", res)
	}
	if cond.Code.Coding[0].Code != "J00" {
		t.Errorf("expected J00, got %s", cond.Code.Coding[0].Code)
	}
}
