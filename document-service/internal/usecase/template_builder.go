package usecase

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
	"github.com/azuvicenna/selaras-services/document-service/templates"
)

var indonesianMonths = map[time.Month]string{
	time.January:   "Januari",
	time.February:  "Februari",
	time.March:     "Maret",
	time.April:     "April",
	time.May:       "Mei",
	time.June:      "Juni",
	time.July:      "Juli",
	time.August:    "Agustus",
	time.September: "September",
	time.October:   "Oktober",
	time.November:  "November",
	time.December:  "Desember",
}

func formatIndonesianDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%d %s %d", t.Day(), indonesianMonths[t.Month()], t.Year())
}

const defaultBaseTypstTemplate = `
#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.5cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123 | Email: rs\@selaras.hospital.id]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Dokumen Resmi RS Selaras Medika - Terverifikasi Sistem Elektronik]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 10pt)

#align(center)[
  #text(12pt, weight: "bold")[{{TITLE}}] \
  #text(9pt)[Nomor: {{DOC_NUMBER}}]
]

#v(1em)

Yang bertanda tangan di bawah ini dokter Rumah Sakit Selaras Medika menerangkan bahwa:

#table(
  columns: (130pt, 10pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}],
  [Nomor Rekam Medis], [:], [{{PATIENT_ID}}],
  [Jenis Kelamin / Tgl Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}],
  [ID Kunjungan / Rawat], [:], [{{ENCOUNTER_ID}}],
)

#v(1em)

{{BODY}}

#v(2em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI KEABSAHAN ELEKTRONIK] \
        #text(6.5pt)[Dokumen ini sah dan dapat diverifikasi melalui:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Pemeriksa, \
      #v(3em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
`

func getEmbeddedTemplateForType(t domain.DocumentType) string {
	var filename string
	switch t {
	// Surat Medis (1-10)
	case domain.DocumentTypeSickLeave:
		filename = "sick_leave.typ"
	case domain.DocumentTypeHealthCertificate:
		filename = "health_certificate.typ"
	case domain.DocumentTypeInpatientStay:
		filename = "inpatient_stay.typ"
	case domain.DocumentTypeReferral:
		filename = "referral.typ"
	case domain.DocumentTypeBirthCertificate:
		filename = "birth_certificate.typ"
	case domain.DocumentTypeDeathCertificate:
		filename = "death_certificate.typ"
	case domain.DocumentTypeInformedConsent:
		filename = "informed_consent.typ"
	case domain.DocumentTypeInformedRefusal:
		filename = "informed_refusal.typ"
	case domain.DocumentTypeMaternityLeave:
		filename = "maternity_leave.typ"
	case domain.DocumentTypeDrugFree:
		filename = "drug_free.typ"

	// Penunjang Medis (11-20)
	case domain.DocumentTypeMedicalRecord:
		filename = "medical_record.typ"
	case domain.DocumentTypeDischargeSummary:
		filename = "discharge_summary.typ"
	case domain.DocumentTypeLabReport:
		filename = "lab_report.typ"
	case domain.DocumentTypeRadiologyReport:
		filename = "radiology_report.typ"
	case domain.DocumentTypeEKGEEGReport:
		filename = "ekg_eeg_report.typ"
	case domain.DocumentTypePathologyReport:
		filename = "pathology_report.typ"
	case domain.DocumentTypePrescription:
		filename = "prescription.typ"
	case domain.DocumentTypeNursingNote:
		filename = "nursing_note.typ"
	case domain.DocumentTypeCarePlan:
		filename = "care_plan.typ"
	case domain.DocumentTypeOdontogram:
		filename = "odontogram.typ"
	}

	if filename != "" {
		if content, err := templates.GetTemplate(filename); err == nil && content != "" {
			return content
		}
	}
	return defaultBaseTypstTemplate
}

func buildTypstSource(doc *domain.MedicalDocument, templateContent, verifyBaseURL string) string {
	tpl := templateContent
	if strings.TrimSpace(tpl) == "" {
		tpl = getEmbeddedTemplateForType(doc.Type)
	}

	issuedDate := ""
	if doc.IssuedAt != nil {
		issuedDate = formatIndonesianDate(*doc.IssuedAt)
	} else {
		issuedDate = formatIndonesianDate(time.Now())
	}

	validUntil := ""
	if doc.ValidUntil != nil {
		validUntil = formatIndonesianDate(*doc.ValidUntil)
	}

	verifyURL := doc.QRCodeURL
	if verifyURL == "" && doc.VerificationToken != "" {
		verifyURL = fmt.Sprintf("%s?token=%s", verifyBaseURL, doc.VerificationToken)
	}

	// Standard substitutions
	replacements := map[string]string{
		"{{TITLE}}":          doc.Title,
		"{{DOC_NUMBER}}":     doc.DocumentNumber,
		"{{PATIENT_ID}}":     doc.PatientID,
		"{{PATIENT_NAME}}":   doc.PatientName,
		"{{PATIENT_GENDER}}": doc.PatientGender,
		"{{PATIENT_DOB}}":    doc.PatientDOB,
		"{{ENCOUNTER_ID}}":   doc.EncounterID,
		"{{DOCTOR_NAME}}":    doc.DoctorName,
		"{{DOCTOR_SIP}}":     doc.DoctorSIP,
		"{{VERIFY_URL}}":     verifyURL,
		"{{ISSUED_DATE}}":    issuedDate,
		"{{VALID_UNTIL}}":    validUntil,
	}

	// Parse JSON payload to support arbitrary custom fields
	var payload map[string]any
	if doc.PayloadJSON != "" {
		if err := json.Unmarshal([]byte(doc.PayloadJSON), &payload); err == nil {
			for k, v := range payload {
				placeholder := fmt.Sprintf("{{%s}}", strings.ToUpper(k))
				replacements[placeholder] = fmt.Sprintf("%v", v)
			}
		}
	}

	// Default body if template has {{BODY}} and payload has notes/diagnosis/summary
	if _, ok := replacements["{{BODY}}"]; !ok {
		var bodyBuilder strings.Builder
		if diagnosis, ok := payload["diagnosis"].(string); ok && diagnosis != "" {
			bodyBuilder.WriteString(fmt.Sprintf("*Diagnosa:* %s\n\n", diagnosis))
		}
		if notes, ok := payload["notes"].(string); ok && notes != "" {
			bodyBuilder.WriteString(fmt.Sprintf("%s\n\n", notes))
		}
		if restDays, ok := payload["rest_days"]; ok {
			bodyBuilder.WriteString(fmt.Sprintf("Diberikan istirahat selama %v hari terhitung sejak tanggal pemeriksaan.\n\n", restDays))
		}
		if bodyBuilder.Len() == 0 {
			bodyBuilder.WriteString("Keterangan klinis telah dicatat dalam sistem rekam medis elektronik.")
		}
		replacements["{{BODY}}"] = bodyBuilder.String()
	}

	res := tpl
	for k, v := range replacements {
		res = strings.ReplaceAll(res, k, v)
	}

	// Clean up any remaining unpopulated template placeholders
	res = remainingPlaceholderRegex.ReplaceAllStringFunc(res, func(match string) string {
		sub := match[2 : len(match)-2]
		if strings.HasSuffix(sub, "_ROWS") || strings.HasSuffix(sub, "_ITEMS") {
			return ""
		}
		return "–"
	})

	return res
}

var remainingPlaceholderRegex = regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`)
