package domain

import "time"

// Partner represents an external enterprise partner (asuransi swasta / korporasi).
type Partner struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Tier       string    `json:"tier"` // basic, premium, enterprise
	APIKey     string    `json:"api_key"`
	SecretKey  string    `json:"-"`
	AllowedIPs []string  `json:"allowed_ips,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

// --- Insurance & Eligibility ---

type InsuranceEligibilityCheckRequest struct {
	PartnerID        string `json:"partner_id"`
	PolicyNumber     string `json:"policy_number"`
	MemberCardNumber string `json:"member_card_number"`
	NIK              string `json:"nik"`
	PatientName      string `json:"patient_name"`
	DateOfBirth      string `json:"date_of_birth"`
	BenefitType      string `json:"benefit_type"` // INPATIENT, OUTPATIENT, DENTAL, OPTICAL, MATERNITY
}

type InsuranceEligibilityCheckResponse struct {
	Eligible              bool    `json:"eligible"`
	MemberName            string  `json:"member_name"`
	CoverageStatus        string  `json:"coverage_status"` // ACTIVE, SUSPENDED, TERMINATED
	ClassLimit            string  `json:"class_limit"`     // VIP, KELAS_1, KELAS_2, KELAS_3
	AnnualLimitRemaining  float64 `json:"annual_limit_remaining"`
	Deductible            float64 `json:"deductible"`
	CoPayPercentage       float64 `json:"copay_percentage"`
	ExpiryDate            string  `json:"expiry_date"`
	Notes                 string  `json:"notes,omitempty"`
}

// --- Pre-Authorization ---

type PreAuthRequest struct {
	PartnerID         string  `json:"partner_id"`
	MemberCardNumber  string  `json:"member_card_number"`
	PatientID         string  `json:"patient_id"`
	AdmissionType     string  `json:"admission_type"` // RAWAT_JALAN, RAWAT_INAP, UGD
	DiagnosaICD10     string  `json:"diagnosa_icd10"`
	ProposedProcedure string  `json:"proposed_procedure"`
	EstimatedCost     float64 `json:"estimated_cost"`
	DoctorName        string  `json:"doctor_name"`
}

type PreAuthResponse struct {
	PreAuthNumber  string  `json:"pre_auth_number"`
	Status         string  `json:"status"` // APPROVED, REJECTED, PENDING_REVIEW
	ApprovedAmount float64 `json:"approved_amount"`
	ValidUntil     string  `json:"valid_until"`
	Remarks        string  `json:"remarks,omitempty"`
}

// --- Claim Ingestion & Tracking ---

type ClaimSubmissionRequest struct {
	PartnerID         string   `json:"partner_id"`
	PreAuthNumber     string   `json:"pre_auth_number,omitempty"`
	PatientID         string   `json:"patient_id"`
	InvoiceID         string   `json:"invoice_id"`
	TotalClaimAmount  float64  `json:"total_claim_amount"`
	DiagnosisCode     string   `json:"diagnosis_code"`
	SupportingDocURLs []string `json:"supporting_doc_urls,omitempty"`
}

type ClaimStatusResponse struct {
	ClaimID             string  `json:"claim_id"`
	ClaimNumber         string  `json:"claim_number"`
	Status              string  `json:"status"` // SUBMITTED, VERIFIED, IN_REVIEW, APPROVED, PAID, REJECTED
	ApprovedAmount      float64 `json:"approved_amount"`
	PaidAt              string  `json:"paid_at,omitempty"`
	SettlementReference string  `json:"settlement_reference,omitempty"`
	RejectionReason     string  `json:"rejection_reason,omitempty"`
}

// --- Corporate Medical Check Up (MCU) ---

type CorporateEmployeeMCU struct {
	EmployeeID  string `json:"employee_id"`
	NIK         string `json:"nik"`
	FullName    string `json:"full_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Phone       string `json:"phone"`
	Department  string `json:"department"`
}

type CorporateMCUBatchRequest struct {
	CorporateID   string                 `json:"corporate_id"`
	CorporateName string                 `json:"corporate_name"`
	PackageCode   string                 `json:"package_code"` // MCU-EXECUTIVE, MCU-STANDARD, MCU-BASIC
	ScheduledDate string                 `json:"scheduled_date"`
	Employees     []CorporateEmployeeMCU `json:"employees"`
}

type CorporateMCUBatchResponse struct {
	BatchID         string `json:"batch_id"`
	TotalEmployees  int    `json:"total_employees"`
	RegisteredCount int    `json:"registered_count"`
	FailedCount     int    `json:"failed_count"`
	BatchStatus     string `json:"batch_status"` // REGISTERED, SCHEDULED, IN_PROGRESS, COMPLETED
}

// --- Webhooks & Event Notification ---

type WebhookRegistration struct {
	WebhookID   string   `json:"webhook_id"`
	PartnerID   string   `json:"partner_id"`
	CallbackURL string   `json:"callback_url"`
	EventTypes  []string `json:"event_types"` // claim.approved, discharge.ready, mcu.completed
	SecretKey   string   `json:"secret_key"`
	IsActive    bool     `json:"is_active"`
}
