package fhir

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/domain"
)

// EncounterParams contains parameters to build an Encounter FHIR resource.
type EncounterParams struct {
	ID               string `json:"id,omitempty"`
	LocalID          string `json:"local_id"`
	PatientIHS       string `json:"patient_ihs"`
	PatientName      string `json:"patient_name,omitempty"`
	PractitionerIHS  string `json:"practitioner_ihs"`
	PractitionerName string `json:"practitioner_name,omitempty"`
	LocationID       string `json:"location_id,omitempty"`
	LocationName     string `json:"location_name,omitempty"`
	Status           string `json:"status"` // arrived, in-progress, finished, cancelled
	ClassCode        string `json:"class_code"` // AMB, IMP, EMER
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time,omitempty"`
}

// ConditionParams contains parameters to build a Condition FHIR resource.
type ConditionParams struct {
	ID           string `json:"id,omitempty"`
	PatientIHS   string `json:"patient_ihs"`
	PatientName  string `json:"patient_name,omitempty"`
	EncounterIHS string `json:"encounter_ihs"`
	ICD10Code    string `json:"icd10_code"`
	Diagnosis    string `json:"diagnosis"`
	RecordedDate string `json:"recorded_date,omitempty"`
}

// ObservationParams contains parameters to build an Observation FHIR resource.
type ObservationParams struct {
	ID              string   `json:"id,omitempty"`
	PatientIHS      string   `json:"patient_ihs"`
	EncounterIHS    string   `json:"encounter_ihs"`
	PractitionerIHS string   `json:"practitioner_ihs,omitempty"`
	LOINCCode       string   `json:"loinc_code"`
	Name            string   `json:"name"`
	Category        string   `json:"category,omitempty"` // vital-signs, laboratory, exam
	ValueQuantity   *float64 `json:"value_quantity,omitempty"`
	Unit            string   `json:"unit,omitempty"`
	UCUMCode        string   `json:"ucum_code,omitempty"`
	ValueString     string   `json:"value_string,omitempty"`
	EffectiveTime   string   `json:"effective_time,omitempty"`
}

// ProcedureParams contains parameters to build a Procedure FHIR resource.
type ProcedureParams struct {
	ID            string `json:"id,omitempty"`
	PatientIHS    string `json:"patient_ihs"`
	EncounterIHS  string `json:"encounter_ihs"`
	ICD9Code      string `json:"icd9_code"`
	ProcedureName string `json:"procedure_name"`
	PerformedTime string `json:"performed_time,omitempty"`
}

// MedicationRequestParams contains parameters to build a MedicationRequest FHIR resource.
type MedicationRequestParams struct {
	ID              string `json:"id,omitempty"`
	PatientIHS      string `json:"patient_ihs"`
	EncounterIHS    string `json:"encounter_ihs"`
	PractitionerIHS string `json:"practitioner_ihs,omitempty"`
	KFACode         string `json:"kfa_code"`
	MedicationName  string `json:"medication_name"`
	Signatura       string `json:"signatura"`
	AuthoredOn      string `json:"authored_on,omitempty"`
}

// DiagnosticReportParams contains parameters to build a DiagnosticReport FHIR resource.
type DiagnosticReportParams struct {
	ID            string `json:"id,omitempty"`
	PatientIHS    string `json:"patient_ihs"`
	EncounterIHS  string `json:"encounter_ihs"`
	LOINCCode     string `json:"loinc_code"`
	TestName      string `json:"test_name"`
	Conclusion    string `json:"conclusion"`
	EffectiveTime string `json:"effective_time,omitempty"`
}

// AllergyIntoleranceParams contains parameters to build an AllergyIntolerance FHIR resource.
type AllergyIntoleranceParams struct {
	ID           string `json:"id,omitempty"`
	PatientIHS   string `json:"patient_ihs"`
	EncounterIHS string `json:"encounter_ihs"`
	SnomedCode   string `json:"snomed_code"`
	AllergenName string `json:"allergen_name"`
	RecordedDate string `json:"recorded_date,omitempty"`
}

// CompositionParams contains parameters to build a Composition FHIR resource.
type CompositionParams struct {
	ID              string `json:"id,omitempty"`
	PatientIHS      string `json:"patient_ihs"`
	EncounterIHS    string `json:"encounter_ihs"`
	PractitionerIHS string `json:"practitioner_ihs"`
	Title           string `json:"title"`
	Date            string `json:"date,omitempty"`
}

// Builder functions

func BuildEncounter(orgID string, p EncounterParams) *EncounterResource {
	status := p.Status
	if status == "" {
		status = "finished"
	}
	classCode := p.ClassCode
	if classCode == "" {
		classCode = "AMB"
	}

	classDisplay := "ambulatory"
	switch classCode {
	case "IMP":
		classDisplay = "inpatient encounter"
	case "EMER":
		classDisplay = "emergency"
	}

	res := &EncounterResource{
		ResourceType: "Encounter",
		ID:           p.ID,
		Status:       status,
		Class: Coding{
			System:  "http://terminology.hl7.org/CodeSystem/v3-ActCode",
			Code:    classCode,
			Display: classDisplay,
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", p.PatientIHS),
			Display:   p.PatientName,
		},
		Period: Period{
			Start: p.StartTime,
			End:   p.EndTime,
		},
	}

	if p.LocalID != "" && orgID != "" {
		res.Identifier = []Identifier{
			{
				System: fmt.Sprintf("http://sys-ids.kemkes.go.id/encounter/%s", orgID),
				Value:  p.LocalID,
			},
		}
	}

	if p.PractitionerIHS != "" {
		res.Participant = []EncounterParticipant{
			{
				Type: []CodeableConcept{
					{
						Coding: []Coding{
							{
								System:  "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
								Code:    "ATND",
								Display: "attender",
							},
						},
					},
				},
				Individual: Reference{
					Reference: fmt.Sprintf("Practitioner/%s", p.PractitionerIHS),
					Display:   p.PractitionerName,
				},
			},
		}
	}

	if p.LocationID != "" {
		res.Location = []EncounterLocation{
			{
				Location: Reference{
					Reference: fmt.Sprintf("Location/%s", p.LocationID),
					Display:   p.LocationName,
				},
			},
		}
	}

	if orgID != "" {
		res.ServiceProvider = Reference{
			Reference: fmt.Sprintf("Organization/%s", orgID),
		}
	}

	return res
}

func BuildCondition(p ConditionParams) *ConditionResource {
	recordedDate := p.RecordedDate
	if recordedDate == "" {
		recordedDate = time.Now().Format(time.RFC3339)
	}

	return &ConditionResource{
		ResourceType: "Condition",
		ID:           p.ID,
		ClinicalStatus: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://terminology.hl7.org/CodeSystem/condition-clinical",
					Code:    "active",
					Display: "Active",
				},
			},
		},
		Category: []CodeableConcept{
			{
				Coding: []Coding{
					{
						System:  "http://terminology.hl7.org/CodeSystem/condition-category",
						Code:    "encounter-diagnosis",
						Display: "Encounter Diagnosis",
					},
				},
			},
		},
		Code: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://hl7.org/fhir/sid/icd-10",
					Code:    p.ICD10Code,
					Display: p.Diagnosis,
				},
			},
			Text: p.Diagnosis,
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", p.PatientIHS),
			Display:   p.PatientName,
		},
		Encounter: Reference{
			Reference: fmt.Sprintf("Encounter/%s", p.EncounterIHS),
		},
		RecordedDate: recordedDate,
	}
}

func BuildObservation(p ObservationParams) *ObservationResource {
	effectiveTime := p.EffectiveTime
	if effectiveTime == "" {
		effectiveTime = time.Now().Format(time.RFC3339)
	}

	cat := p.Category
	if cat == "" {
		cat = "vital-signs"
	}

	res := &ObservationResource{
		ResourceType: "Observation",
		ID:           p.ID,
		Status:       "final",
		Category: []CodeableConcept{
			{
				Coding: []Coding{
					{
						System:  "http://terminology.hl7.org/CodeSystem/observation-category",
						Code:    cat,
						Display: strings.Title(cat),
					},
				},
			},
		},
		Code: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://loinc.org",
					Code:    p.LOINCCode,
					Display: p.Name,
				},
			},
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", p.PatientIHS),
		},
		Encounter: Reference{
			Reference: fmt.Sprintf("Encounter/%s", p.EncounterIHS),
		},
		EffectiveDateTime: effectiveTime,
		Issued:            effectiveTime,
	}

	if p.PractitionerIHS != "" {
		res.Performer = []Reference{
			{
				Reference: fmt.Sprintf("Practitioner/%s", p.PractitionerIHS),
			},
		}
	}

	if p.ValueQuantity != nil {
		ucum := p.UCUMCode
		if ucum == "" {
			ucum = p.Unit
		}
		res.ValueQuantity = &Quantity{
			Value:  *p.ValueQuantity,
			Unit:   p.Unit,
			System: "http://unitsofmeasure.org",
			Code:   ucum,
		}
	} else if p.ValueString != "" {
		res.ValueString = p.ValueString
	}

	return res
}

func BuildProcedure(p ProcedureParams) *ProcedureResource {
	performedTime := p.PerformedTime
	if performedTime == "" {
		performedTime = time.Now().Format(time.RFC3339)
	}

	return &ProcedureResource{
		ResourceType: "Procedure",
		ID:           p.ID,
		Status:       "completed",
		Category: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://snomed.info/sct",
					Code:    "103693007",
					Display: "Diagnostic procedure",
				},
			},
		},
		Code: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://hl7.org/fhir/sid/icd-9-cm",
					Code:    p.ICD9Code,
					Display: p.ProcedureName,
				},
			},
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", p.PatientIHS),
		},
		Encounter: Reference{
			Reference: fmt.Sprintf("Encounter/%s", p.EncounterIHS),
		},
		PerformedDateTime: performedTime,
	}
}

func BuildMedicationRequest(p MedicationRequestParams) *MedicationRequestResource {
	authoredOn := p.AuthoredOn
	if authoredOn == "" {
		authoredOn = time.Now().Format(time.RFC3339)
	}

	res := &MedicationRequestResource{
		ResourceType: "MedicationRequest",
		ID:           p.ID,
		Status:       "active",
		Intent:       "order",
		Category: []CodeableConcept{
			{
				Coding: []Coding{
					{
						System:  "http://terminology.hl7.org/CodeSystem/medicationrequest-category",
						Code:    "outpatient",
						Display: "Outpatient",
					},
				},
			},
		},
		MedicationCodeableConcept: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://sys-ids.kemkes.go.id/kfa",
					Code:    p.KFACode,
					Display: p.MedicationName,
				},
			},
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", p.PatientIHS),
		},
		Encounter: Reference{
			Reference: fmt.Sprintf("Encounter/%s", p.EncounterIHS),
		},
		AuthoredOn: authoredOn,
	}

	if p.PractitionerIHS != "" {
		res.Requester = Reference{
			Reference: fmt.Sprintf("Practitioner/%s", p.PractitionerIHS),
		}
	}

	if p.Signatura != "" {
		res.DosageInstruction = []DosageInstruction{
			{
				Text: p.Signatura,
			},
		}
	}

	return res
}

func BuildDiagnosticReport(p DiagnosticReportParams) *DiagnosticReportResource {
	effectiveTime := p.EffectiveTime
	if effectiveTime == "" {
		effectiveTime = time.Now().Format(time.RFC3339)
	}

	return &DiagnosticReportResource{
		ResourceType: "DiagnosticReport",
		ID:           p.ID,
		Status:       "final",
		Category: []CodeableConcept{
			{
				Coding: []Coding{
					{
						System:  "http://terminology.hl7.org/CodeSystem/v2-0074",
						Code:    "LAB",
						Display: "Laboratory",
					},
				},
			},
		},
		Code: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://loinc.org",
					Code:    p.LOINCCode,
					Display: p.TestName,
				},
			},
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", p.PatientIHS),
		},
		Encounter: Reference{
			Reference: fmt.Sprintf("Encounter/%s", p.EncounterIHS),
		},
		EffectiveDateTime: effectiveTime,
		Issued:            effectiveTime,
		Conclusion:        p.Conclusion,
	}
}

func BuildAllergyIntolerance(p AllergyIntoleranceParams) *AllergyIntoleranceResource {
	recordedDate := p.RecordedDate
	if recordedDate == "" {
		recordedDate = time.Now().Format(time.RFC3339)
	}

	return &AllergyIntoleranceResource{
		ResourceType: "AllergyIntolerance",
		ID:           p.ID,
		ClinicalStatus: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical",
					Code:    "active",
					Display: "Active",
				},
			},
		},
		VerificationStatus: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification",
					Code:    "confirmed",
					Display: "Confirmed",
				},
			},
		},
		Code: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://snomed.info/sct",
					Code:    p.SnomedCode,
					Display: p.AllergenName,
				},
			},
		},
		Patient: Reference{
			Reference: fmt.Sprintf("Patient/%s", p.PatientIHS),
		},
		Encounter: Reference{
			Reference: fmt.Sprintf("Encounter/%s", p.EncounterIHS),
		},
		RecordedDate: recordedDate,
	}
}

func BuildComposition(orgID string, p CompositionParams) *CompositionResource {
	date := p.Date
	if date == "" {
		date = time.Now().Format(time.RFC3339)
	}

	res := &CompositionResource{
		ResourceType: "Composition",
		ID:           p.ID,
		Status:       "final",
		Type: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://loinc.org",
					Code:    "11506-3",
					Display: "Progress note",
				},
			},
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", p.PatientIHS),
		},
		Encounter: Reference{
			Reference: fmt.Sprintf("Encounter/%s", p.EncounterIHS),
		},
		Date:  date,
		Title: p.Title,
	}

	if p.PractitionerIHS != "" {
		res.Author = []Reference{
			{
				Reference: fmt.Sprintf("Practitioner/%s", p.PractitionerIHS),
			},
		}
	}

	if orgID != "" {
		res.Custodian = Reference{
			Reference: fmt.Sprintf("Organization/%s", orgID),
		}
	}

	return res
}

// PrepareFHIRPayload inspects raw JSON and transforms it into a FHIR R4 structure if needed.
func PrepareFHIRPayload(orgID, resType string, raw json.RawMessage) (any, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty raw payload")
	}

	// 1. Check if raw JSON already has "resourceType"
	var probe struct {
		ResourceType string `json:"resourceType"`
	}
	if err := json.Unmarshal(raw, &probe); err == nil && probe.ResourceType != "" {
		// Already a FHIR JSON structure
		return raw, nil
	}

	fhirType, err := domain.ToFHIRResourceType(resType)
	if err != nil {
		return nil, err
	}

	switch fhirType {
	case "Encounter":
		var p EncounterParams
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal encounter params: %w", err)
		}
		return BuildEncounter(orgID, p), nil

	case "Condition":
		var p ConditionParams
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal condition params: %w", err)
		}
		return BuildCondition(p), nil

	case "Observation":
		var p ObservationParams
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal observation params: %w", err)
		}
		return BuildObservation(p), nil

	case "Procedure":
		var p ProcedureParams
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal procedure params: %w", err)
		}
		return BuildProcedure(p), nil

	case "MedicationRequest":
		var p MedicationRequestParams
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal medication params: %w", err)
		}
		return BuildMedicationRequest(p), nil

	case "DiagnosticReport":
		var p DiagnosticReportParams
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal diagnostic report params: %w", err)
		}
		return BuildDiagnosticReport(p), nil

	case "AllergyIntolerance":
		var p AllergyIntoleranceParams
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal allergy params: %w", err)
		}
		return BuildAllergyIntolerance(p), nil

	case "Composition":
		var p CompositionParams
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal composition params: %w", err)
		}
		return BuildComposition(orgID, p), nil

	default:
		return nil, fmt.Errorf("unsupported resource builder: %s", fhirType)
	}
}
