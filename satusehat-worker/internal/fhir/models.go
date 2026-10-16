package fhir

// Common FHIR Data Types
type Coding struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}

type Identifier struct {
	Use    string          `json:"use,omitempty"`
	System string          `json:"system,omitempty"`
	Value  string          `json:"value,omitempty"`
	Type   CodeableConcept `json:"type,omitempty"`
}

type Reference struct {
	Reference string `json:"reference,omitempty"`
	Display   string `json:"display,omitempty"`
}

type Period struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type Quantity struct {
	Value  float64 `json:"value,omitempty"`
	Unit   string  `json:"unit,omitempty"`
	System string  `json:"system,omitempty"`
	Code   string  `json:"code,omitempty"`
}

// Participant for Encounter
type EncounterParticipant struct {
	Type       []CodeableConcept `json:"type,omitempty"`
	Individual Reference         `json:"individual,omitempty"`
}

// Location reference for Encounter
type EncounterLocation struct {
	Location Reference `json:"location,omitempty"`
}

// 1. Encounter FHIR R4 Resource
type EncounterResource struct {
	ResourceType    string                 `json:"resourceType"`
	ID              string                 `json:"id,omitempty"`
	Identifier      []Identifier           `json:"identifier,omitempty"`
	Status          string                 `json:"status"`
	Class           Coding                 `json:"class"`
	Subject         Reference              `json:"subject"`
	Participant     []EncounterParticipant `json:"participant,omitempty"`
	Period          Period                 `json:"period"`
	Location        []EncounterLocation    `json:"location,omitempty"`
	ServiceProvider Reference              `json:"serviceProvider,omitempty"`
}

// 2. Condition (Diagnosis) FHIR R4 Resource
type ConditionResource struct {
	ResourceType   string            `json:"resourceType"`
	ID             string            `json:"id,omitempty"`
	ClinicalStatus CodeableConcept   `json:"clinicalStatus"`
	Category       []CodeableConcept `json:"category,omitempty"`
	Code           CodeableConcept   `json:"code"`
	Subject        Reference         `json:"subject"`
	Encounter      Reference         `json:"encounter"`
	RecordedDate   string            `json:"recordedDate,omitempty"`
}

// 3. Observation FHIR R4 Resource
type ObservationResource struct {
	ResourceType      string            `json:"resourceType"`
	ID                string            `json:"id,omitempty"`
	Status            string            `json:"status"`
	Category          []CodeableConcept `json:"category,omitempty"`
	Code              CodeableConcept   `json:"code"`
	Subject           Reference         `json:"subject"`
	Encounter         Reference         `json:"encounter"`
	EffectiveDateTime string            `json:"effectiveDateTime,omitempty"`
	Issued            string            `json:"issued,omitempty"`
	Performer         []Reference       `json:"performer,omitempty"`
	ValueQuantity     *Quantity         `json:"valueQuantity,omitempty"`
	ValueString       string            `json:"valueString,omitempty"`
}

// 4. Procedure FHIR R4 Resource
type ProcedureResource struct {
	ResourceType      string          `json:"resourceType"`
	ID                string          `json:"id,omitempty"`
	Status            string          `json:"status"`
	Category          CodeableConcept `json:"category,omitempty"`
	Code              CodeableConcept `json:"code"`
	Subject           Reference       `json:"subject"`
	Encounter         Reference       `json:"encounter"`
	PerformedDateTime string          `json:"performedDateTime,omitempty"`
}

// 5. MedicationRequest FHIR R4 Resource
type DosageInstruction struct {
	Text string `json:"text,omitempty"`
}

type MedicationRequestResource struct {
	ResourceType              string              `json:"resourceType"`
	ID                        string              `json:"id,omitempty"`
	Status                    string              `json:"status"`
	Intent                    string              `json:"intent"`
	Category                  []CodeableConcept   `json:"category,omitempty"`
	MedicationCodeableConcept CodeableConcept     `json:"medicationCodeableConcept"`
	Subject                   Reference           `json:"subject"`
	Encounter                 Reference           `json:"encounter"`
	AuthoredOn                string              `json:"authoredOn,omitempty"`
	Requester                 Reference           `json:"requester,omitempty"`
	DosageInstruction         []DosageInstruction `json:"dosageInstruction,omitempty"`
}

// 6. DiagnosticReport FHIR R4 Resource
type DiagnosticReportResource struct {
	ResourceType      string            `json:"resourceType"`
	ID                string            `json:"id,omitempty"`
	Status            string            `json:"status"`
	Category          []CodeableConcept `json:"category,omitempty"`
	Code              CodeableConcept   `json:"code"`
	Subject           Reference         `json:"subject"`
	Encounter         Reference         `json:"encounter"`
	EffectiveDateTime string            `json:"effectiveDateTime,omitempty"`
	Issued            string            `json:"issued,omitempty"`
	Conclusion        string            `json:"conclusion,omitempty"`
}

// 7. AllergyIntolerance FHIR R4 Resource
type AllergyIntoleranceResource struct {
	ResourceType       string          `json:"resourceType"`
	ID                 string          `json:"id,omitempty"`
	ClinicalStatus     CodeableConcept `json:"clinicalStatus"`
	VerificationStatus CodeableConcept `json:"verificationStatus"`
	Code               CodeableConcept `json:"code"`
	Patient            Reference       `json:"patient"`
	Encounter          Reference       `json:"encounter"`
	RecordedDate       string          `json:"recordedDate,omitempty"`
}

// 8. Composition (ClinicalNote) FHIR R4 Resource
type CompositionResource struct {
	ResourceType string          `json:"resourceType"`
	ID           string          `json:"id,omitempty"`
	Status       string          `json:"status"`
	Type         CodeableConcept `json:"type"`
	Subject      Reference       `json:"subject"`
	Encounter    Reference       `json:"encounter"`
	Date         string          `json:"date"`
	Author       []Reference     `json:"author"`
	Title        string          `json:"title"`
	Custodian    Reference       `json:"custodian,omitempty"`
}
