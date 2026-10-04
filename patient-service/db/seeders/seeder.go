package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	patientv1 "github.com/azuvicenna/selaras-services/patient-service/proto/patient/v1"
)

type SeedPatientWrapper struct {
	Patient          json.RawMessage   `json:"patient"`
	Allergies        []json.RawMessage `json:"allergies"`
	Documents        []json.RawMessage `json:"documents"`
	EmergencyContact json.RawMessage   `json:"emergencyContact"`
	Families         []json.RawMessage `json:"families"`
	Insurances       []json.RawMessage `json:"insurances"`
}

type PatientAggregate struct {
	Patient          *patientv1.Patient
	Allergies        []*patientv1.PatientAllergy
	Documents        []*patientv1.PatientDocument
	EmergencyContact *patientv1.PatientEmergencyContact
	Families         []*patientv1.PatientFamily
	Insurances       []*patientv1.PatientInsurance
}

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
		now := timestamppb.Now()

		// 1. Patient
		agg.Patient = &patientv1.Patient{}
		if err := protojson.Unmarshal(item.Patient, agg.Patient); err != nil {
			return nil, fmt.Errorf("failed to unmarshal patient: %w", err)
		}
		patientID := ulid.Make().String()
		agg.Patient.Id = patientID
		agg.Patient.CreatedAt = now
		agg.Patient.UpdatedAt = now

		// 2. Allergies
		for _, raw := range item.Allergies {
			allergy := &patientv1.PatientAllergy{}
			if err := protojson.Unmarshal(raw, allergy); err != nil {
				return nil, fmt.Errorf("failed to unmarshal allergy: %w", err)
			}
			allergy.Id = ulid.Make().String()
			allergy.PatientId = patientID
			allergy.CreatedAt = now
			allergy.UpdatedAt = now
			agg.Allergies = append(agg.Allergies, allergy)
		}

		// 3. Documents
		for _, raw := range item.Documents {
			doc := &patientv1.PatientDocument{}
			if err := protojson.Unmarshal(raw, doc); err != nil {
				return nil, fmt.Errorf("failed to unmarshal document: %w", err)
			}
			doc.Id = ulid.Make().String()
			doc.PatientId = patientID
			doc.CreatedAt = now
			doc.UpdatedAt = now
			agg.Documents = append(agg.Documents, doc)
		}

		// 4. Emergency Contact
		if len(item.EmergencyContact) > 0 && string(item.EmergencyContact) != "null" {
			contact := &patientv1.PatientEmergencyContact{}
			if err := protojson.Unmarshal(item.EmergencyContact, contact); err != nil {
				return nil, fmt.Errorf("failed to unmarshal emergency contact: %w", err)
			}
			contact.Id = ulid.Make().String()
			contact.PatientId = patientID
			contact.CreatedAt = now
			contact.UpdatedAt = now
			agg.EmergencyContact = contact
		}

		// 5. Families
		for _, raw := range item.Families {
			family := &patientv1.PatientFamily{}
			if err := protojson.Unmarshal(raw, family); err != nil {
				return nil, fmt.Errorf("failed to unmarshal family: %w", err)
			}
			family.Id = ulid.Make().String()
			family.PatientId = patientID
			family.CreatedAt = now
			family.UpdatedAt = now
			agg.Families = append(agg.Families, family)
		}

		// 6. Insurances
		for _, raw := range item.Insurances {
			insurance := &patientv1.PatientInsurance{}
			if err := protojson.Unmarshal(raw, insurance); err != nil {
				return nil, fmt.Errorf("failed to unmarshal insurance: %w", err)
			}
			insurance.Id = ulid.Make().String()
			insurance.PatientId = patientID
			insurance.CreatedAt = now
			insurance.UpdatedAt = now
			agg.Insurances = append(agg.Insurances, insurance)
		}

		aggregates = append(aggregates, &agg)
	}

	return aggregates, nil
}