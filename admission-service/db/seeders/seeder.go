package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedSchedule struct {
	PractitionerID string `json:"practitioner_id"`
	DepartmentID   string `json:"department_id"`
	DayOfWeek      int32  `json:"day_of_week"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Quota          int32  `json:"quota"`
	IsActive       bool   `json:"is_active"`
}

type SeedAppointment struct {
	BookingCode     string `json:"booking_code"`
	PatientID       string `json:"patient_id"`
	PractitionerID  string `json:"practitioner_id"`
	DepartmentID    string `json:"department_id"`
	ScheduleID      string `json:"schedule_id"`
	AppointmentDate string `json:"appointment_date"`
	QueueNumber     int32  `json:"queue_number"`
	Status          int32  `json:"status"`
	Notes           string `json:"notes"`
}

type SeedAdmission struct {
	AdmissionNumber string `json:"admission_number"`
	PatientID       string `json:"patient_id"`
	AppointmentID   string `json:"appointment_id"`
	AdmissionType   int32  `json:"admission_type"`
	Status          int32  `json:"status"`
	DepartmentID    string `json:"department_id"`
	PractitionerID  string `json:"practitioner_id"`
	WardID          string `json:"ward_id"`
	BedID           string `json:"bed_id"`
	Notes           string `json:"notes"`
}

type SeedData struct {
	Schedules    []SeedSchedule    `json:"schedules"`
	Appointments []SeedAppointment `json:"appointments"`
	Admissions   []SeedAdmission   `json:"admissions"`
}

func loadSeedData(path string) (*SeedData, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	var data SeedData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("unmarshal seed json: %w", err)
	}

	return &data, nil
}
