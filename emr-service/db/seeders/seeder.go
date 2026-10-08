package main

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/protobuf/encoding/protojson"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
)

type SeedEncounterWrapper struct {
	Encounter          json.RawMessage   `json:"encounter"`
	Allergies          []json.RawMessage `json:"allergies"`
	ClinicalNotes      []json.RawMessage `json:"clinicalNotes"`
	Conditions         []json.RawMessage `json:"conditions"`
	Observations       []json.RawMessage `json:"observations"`
	Procedures         []json.RawMessage `json:"procedures"`
	MedicationRequests []json.RawMessage `json:"medicationRequests"`
	DiagnosticReports  []json.RawMessage `json:"diagnosticReports"`
}

type EncounterAggregate struct {
	Encounter          *emrv1.Encounter
	Allergies          []*emrv1.Allergy
	ClinicalNotes      []*emrv1.ClinicalNote
	Conditions         []*emrv1.Condition
	Observations       []*emrv1.Observation
	Procedures         []*emrv1.Procedure
	MedicationRequests []*emrv1.MedicationRequest
	DiagnosticReports  []*emrv1.DiagnosticReport
}

// ID, encounter_id, dan timestamp tidak wajib diisi di JSON; service yang membuatnya saat seeding.
func loadSeedData(filePath string) ([]*EncounterAggregate, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var rawItems []SeedEncounterWrapper
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON wrapper: %w", err)
	}

	// protojson menerima enum berupa string ("ENCOUNTER_CLASS_AMBULATORY") maupun angka (1)
	var aggregates []*EncounterAggregate

	for _, item := range rawItems {
		var agg EncounterAggregate

		agg.Encounter = &emrv1.Encounter{}
		if err := protojson.Unmarshal(item.Encounter, agg.Encounter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal encounter: %w", err)
		}

		for _, raw := range item.Allergies {
			allergy := &emrv1.Allergy{}
			if err := protojson.Unmarshal(raw, allergy); err != nil {
				return nil, fmt.Errorf("failed to unmarshal allergy: %w", err)
			}
			agg.Allergies = append(agg.Allergies, allergy)
		}

		for _, raw := range item.ClinicalNotes {
			note := &emrv1.ClinicalNote{}
			if err := protojson.Unmarshal(raw, note); err != nil {
				return nil, fmt.Errorf("failed to unmarshal clinical note: %w", err)
			}
			agg.ClinicalNotes = append(agg.ClinicalNotes, note)
		}

		for _, raw := range item.Conditions {
			cond := &emrv1.Condition{}
			if err := protojson.Unmarshal(raw, cond); err != nil {
				return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
			}
			agg.Conditions = append(agg.Conditions, cond)
		}

		for _, raw := range item.Observations {
			obs := &emrv1.Observation{}
			if err := protojson.Unmarshal(raw, obs); err != nil {
				return nil, fmt.Errorf("failed to unmarshal observation: %w", err)
			}
			agg.Observations = append(agg.Observations, obs)
		}

		for _, raw := range item.Procedures {
			proc := &emrv1.Procedure{}
			if err := protojson.Unmarshal(raw, proc); err != nil {
				return nil, fmt.Errorf("failed to unmarshal procedure: %w", err)
			}
			agg.Procedures = append(agg.Procedures, proc)
		}

		for _, raw := range item.MedicationRequests {
			med := &emrv1.MedicationRequest{}
			if err := protojson.Unmarshal(raw, med); err != nil {
				return nil, fmt.Errorf("failed to unmarshal medication request: %w", err)
			}
			agg.MedicationRequests = append(agg.MedicationRequests, med)
		}

		for _, raw := range item.DiagnosticReports {
			report := &emrv1.DiagnosticReport{}
			if err := protojson.Unmarshal(raw, report); err != nil {
				return nil, fmt.Errorf("failed to unmarshal diagnostic report: %w", err)
			}
			agg.DiagnosticReports = append(agg.DiagnosticReports, report)
		}

		aggregates = append(aggregates, &agg)
	}

	return aggregates, nil
}
