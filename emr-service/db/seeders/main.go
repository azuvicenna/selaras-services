package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50052", "alamat gRPC emr-service")
	file := flag.String("file", "db/seeders/emr.json", "file JSON data seed")
	flag.Parse()

	if err := run(*addr, *file); err != nil {
		slog.Error("seeding failed", "err", err)
		os.Exit(1)
	}
}

func run(addr, file string) error {
	aggregates, err := loadSeedData(file)
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to create gRPC client: %w", err)
	}
	defer func() { _ = conn.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client := emrv1.NewEmrServiceClient(conn)
	for _, agg := range aggregates {
		if err := seedEncounter(ctx, client, agg); err != nil {
			return fmt.Errorf("seed encounter for patient %q: %w", agg.Encounter.GetPatientId(), err)
		}
	}

	slog.Info("seeding finished", "encounters", len(aggregates))
	return nil
}

func seedEncounter(ctx context.Context, c emrv1.EmrServiceClient, agg *EncounterAggregate) error {
	res, err := c.StartEncounter(ctx, &emrv1.StartEncounterRequest{Encounter: agg.Encounter})
	if status.Code(err) == codes.AlreadyExists {
		slog.Info("skipped, encounter already exists", "satusehat_id", agg.Encounter.GetSatusehatId())
		return nil
	}
	if err != nil {
		return err
	}

	encounterID := res.GetEncounter().GetId()
	patientID := agg.Encounter.GetPatientId()
	practitionerID := agg.Encounter.GetPractitionerId()

	for _, allergy := range agg.Allergies {
		allergy.PatientId = patientID
		allergy.EncounterId = encounterID
		if allergy.PractitionerId == "" {
			allergy.PractitionerId = practitionerID
		}
		if _, err := c.CreatePatientAllergy(ctx, &emrv1.CreatePatientAllergyRequest{Allergy: allergy}); err != nil {
			return fmt.Errorf("allergy: %w", err)
		}
	}

	var observationIDs []string
	for _, obs := range agg.Observations {
		obs.PatientId = patientID
		obs.EncounterId = encounterID
		if obs.PractitionerId == "" {
			obs.PractitionerId = practitionerID
		}
		obsRes, err := c.CreateObservation(ctx, &emrv1.CreateObservationRequest{Observation: obs})
		if err != nil {
			return fmt.Errorf("observation: %w", err)
		}
		if id := obsRes.GetObservation().GetId(); id != "" {
			observationIDs = append(observationIDs, id)
		}
	}

	var primaryNoteID string
	for _, note := range agg.ClinicalNotes {
		note.PatientId = patientID
		note.EncounterId = encounterID
		if note.PractitionerId == "" {
			note.PractitionerId = practitionerID
		}
		if len(note.ObservationIds) == 0 && len(observationIDs) > 0 {
			note.ObservationIds = observationIDs
		}
		noteRes, err := c.CreateClinicalNote(ctx, &emrv1.CreateClinicalNoteRequest{ClinicalNote: note})
		if err != nil {
			return fmt.Errorf("clinical note: %w", err)
		}
		if primaryNoteID == "" {
			primaryNoteID = noteRes.GetClinicalNote().GetId()
		}
	}

	var primaryConditionID string
	for _, cond := range agg.Conditions {
		cond.PatientId = patientID
		cond.EncounterId = encounterID
		if cond.PractitionerId == "" {
			cond.PractitionerId = practitionerID
		}
		if cond.ClinicalNoteId == "" {
			cond.ClinicalNoteId = primaryNoteID
		}
		condRes, err := c.AddCondition(ctx, &emrv1.AddConditionRequest{Condition: cond})
		if err != nil {
			return fmt.Errorf("condition: %w", err)
		}
		if primaryConditionID == "" || cond.GetIsPrimary() {
			primaryConditionID = condRes.GetCondition().GetId()
		}
	}

	for _, proc := range agg.Procedures {
		proc.PatientId = patientID
		proc.EncounterId = encounterID
		if proc.PractitionerId == "" {
			proc.PractitionerId = practitionerID
		}
		if proc.ReasonConditionId == "" {
			proc.ReasonConditionId = primaryConditionID
		}
		if proc.ClinicalNoteId == "" {
			proc.ClinicalNoteId = primaryNoteID
		}
		if _, err := c.RecordProcedure(ctx, &emrv1.RecordProcedureRequest{Procedure: proc}); err != nil {
			return fmt.Errorf("procedure: %w", err)
		}
	}

	for _, med := range agg.MedicationRequests {
		med.PatientId = patientID
		med.EncounterId = encounterID
		if med.PractitionerId == "" {
			med.PractitionerId = practitionerID
		}
		if med.ConditionId == "" {
			med.ConditionId = primaryConditionID
		}
		if _, err := c.CreateMedicationRequest(ctx, &emrv1.CreateMedicationRequestRequest{MedicationRequest: med}); err != nil {
			return fmt.Errorf("medication request: %w", err)
		}
	}

	for _, report := range agg.DiagnosticReports {
		report.PatientId = patientID
		report.EncounterId = encounterID
		if report.RequesterId == "" {
			report.RequesterId = practitionerID
		}
		if _, err := c.CreateDiagnosticReport(ctx, &emrv1.CreateDiagnosticReportRequest{DiagnosticReport: report}); err != nil {
			return fmt.Errorf("diagnostic report: %w", err)
		}
	}

	if agg.Encounter.GetStatus() == emrv1.EncounterStatus_ENCOUNTER_STATUS_FINISHED {
		if _, err := c.FinishEncounter(ctx, &emrv1.FinishEncounterRequest{
			Id:                   encounterID,
			DischargeDisposition: agg.Encounter.GetDischargeDisposition(),
		}); err != nil {
			return fmt.Errorf("finish encounter: %w", err)
		}
	}

	slog.Info("seeded encounter", "patient_id", patientID, "encounter_id", encounterID)
	return nil
}
