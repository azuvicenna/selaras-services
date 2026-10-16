package domain

type MobilePatientProfile struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	NIK                 string `json:"nik"`
	MedicalRecordNumber string `json:"medical_record_number"`
	Phone               string `json:"phone"`
	BirthDate           string `json:"birth_date"`
	Gender              string `json:"gender"`
	InsuranceName       string `json:"insurance_name,omitempty"`
	InsuranceCardNumber string `json:"insurance_card_number,omitempty"`
}

type MobileActiveQueue struct {
	TicketID            string `json:"ticket_id"`
	TicketNumber        string `json:"ticket_number"`
	PoliName            string `json:"poli_name"`
	Destination         string `json:"destination"`
	CurrentCalledNumber string `json:"current_called_number"`
	RemainingAhead      int    `json:"remaining_ahead"`
	EstimatedWaitTime   string `json:"estimated_wait_time"`
	Status              string `json:"status"`
}

type MobileUpcomingAppointment struct {
	AppointmentID    string `json:"appointment_id"`
	DoctorName       string `json:"doctor_name"`
	PoliName         string `json:"poli_name"`
	ScheduleTime     string `json:"schedule_time"`
	Status           string `json:"status"`
	QueueNumber      string `json:"queue_number,omitempty"`
}

type MobileHomeDashboard struct {
	Profile             *MobilePatientProfile      `json:"profile,omitempty"`
	ActiveQueue         *MobileActiveQueue         `json:"active_queue,omitempty"`
	UpcomingAppointment *MobileUpcomingAppointment `json:"upcoming_appointment,omitempty"`
	UnreadAlertsCount   int                        `json:"unread_alerts_count"`
	AvailableBedsCount  int                        `json:"available_beds_count"`
}

type MobileDoctorSchedule struct {
	PractitionerID string `json:"practitioner_id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
	PoliName       string `json:"poli_name"`
	DayOfWeek      string `json:"day_of_week"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	QuotaAvailable int    `json:"quota_available"`
}

type MobileBookingRequest struct {
	PatientID      string `json:"patient_id"`
	PractitionerID string `json:"practitioner_id"`
	PoliID         string `json:"poli_id"`
	BookingDate    string `json:"booking_date"`
	PaymentMethod  string `json:"payment_method"` // BPJS, UMUM, ASURANSI
	ReferralNumber string `json:"referral_number,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

type MobileBookingResponse struct {
	ReservationID string `json:"reservation_id"`
	QueueNumber   string `json:"queue_number"`
	EstimatedTime string `json:"estimated_time"`
	PoliName      string `json:"poli_name"`
	DoctorName    string `json:"doctor_name"`
	BookingDate   string `json:"booking_date"`
}

type MobileMedicalRecordItem struct {
	EncounterID     string   `json:"encounter_id"`
	Date            string   `json:"date"`
	DoctorName      string   `json:"doctor_name"`
	PoliName        string   `json:"poli_name"`
	ChiefComplaint  string   `json:"chief_complaint"`
	Diagnosis       string   `json:"diagnosis"`
	Prescriptions   []string `json:"prescriptions,omitempty"`
	HasLabResult    bool     `json:"has_lab_result"`
	HasRadiology    bool     `json:"has_radiology"`
	LabReportURL    string   `json:"lab_report_url,omitempty"`
	RadiologyReportURL string `json:"radiology_report_url,omitempty"`
}
