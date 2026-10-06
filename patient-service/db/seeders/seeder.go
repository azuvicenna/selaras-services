package main

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/protobuf/encoding/protojson"

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
)

// Dokumen tidak ikut di-seed: UploadPatientDocument butuh isi file dan storage.
// Field "documents" di JSON diabaikan.
type SeedPatientWrapper struct {
	Patient          json.RawMessage   `json:"patient"`
	Allergies        []json.RawMessage `json:"allergies"`
	EmergencyContact json.RawMessage   `json:"emergencyContact"`
	Families         []json.RawMessage `json:"families"`
	Insurances       []json.RawMessage `json:"insurances"`
}

type PatientAggregate struct {
	Patient          *patientv1.Patient
	Allergies        []*patientv1.PatientAllergy
	EmergencyContact *patientv1.PatientEmergencyContact
	Families         []*patientv1.PatientFamily
	Insurances       []*patientv1.PatientInsurance
}

// ID, patient_id, dan timestamp tidak diisi di sini; service yang membuatnya.
func loadSeedData(filePath string) ([]*PatientAggregate, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var rawItems []SeedPatientWrapper
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON wrapper: %w", err)
	}

	// protojson menerima enum berupa string ("GENDER_FEMALE") maupun angka (2)
	var aggregates []*PatientAggregate

	for _, item := range rawItems {
		var agg PatientAggregate

		agg.Patient = &patientv1.Patient{}
		if err := protojson.Unmarshal(item.Patient, agg.Patient); err != nil {
			return nil, fmt.Errorf("failed to unmarshal patient: %w", err)
		}

		for _, raw := range item.Allergies {
			allergy := &patientv1.PatientAllergy{}
			if err := protojson.Unmarshal(raw, allergy); err != nil {
				return nil, fmt.Errorf("failed to unmarshal allergy: %w", err)
			}
			agg.Allergies = append(agg.Allergies, allergy)
		}

		if len(item.EmergencyContact) > 0 && string(item.EmergencyContact) != "null" {
			contact := &patientv1.PatientEmergencyContact{}
			if err := protojson.Unmarshal(item.EmergencyContact, contact); err != nil {
				return nil, fmt.Errorf("failed to unmarshal emergency contact: %w", err)
			}
			agg.EmergencyContact = contact
		}

		for _, raw := range item.Families {
			family := &patientv1.PatientFamily{}
			if err := protojson.Unmarshal(raw, family); err != nil {
				return nil, fmt.Errorf("failed to unmarshal family: %w", err)
			}
			agg.Families = append(agg.Families, family)
		}

		for _, raw := range item.Insurances {
			insurance := &patientv1.PatientInsurance{}
			if err := protojson.Unmarshal(raw, insurance); err != nil {
				return nil, fmt.Errorf("failed to unmarshal insurance: %w", err)
			}
			agg.Insurances = append(agg.Insurances, insurance)
		}

		aggregates = append(aggregates, &agg)
	}

	return aggregates, nil
}
