package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
)

func getPrefixForType(t domain.DocumentType) string {
	switch t {
	case domain.DocumentTypeSickLeave:
		return "SKS"
	case domain.DocumentTypeHealthCertificate:
		return "SKD"
	case domain.DocumentTypeInpatientStay:
		return "SKRI"
	case domain.DocumentTypeReferral:
		return "RUJ"
	case domain.DocumentTypeBirthCertificate:
		return "SKL"
	case domain.DocumentTypeDeathCertificate:
		return "SKK"
	case domain.DocumentTypeInformedConsent:
		return "IC"
	case domain.DocumentTypeInformedRefusal:
		return "IR"
	case domain.DocumentTypeMaternityLeave:
		return "SKC"
	case domain.DocumentTypeDrugFree:
		return "SKBN"
	case domain.DocumentTypeMedicalRecord:
		return "RM"
	case domain.DocumentTypeDischargeSummary:
		return "RES"
	case domain.DocumentTypeLabReport:
		return "LAB"
	case domain.DocumentTypeRadiologyReport:
		return "RAD"
	case domain.DocumentTypeEKGEEGReport:
		return "EKG"
	case domain.DocumentTypePathologyReport:
		return "PA"
	case domain.DocumentTypePrescription:
		return "RSP"
	case domain.DocumentTypeNursingNote:
		return "CPPT"
	case domain.DocumentTypeCarePlan:
		return "CP"
	case domain.DocumentTypeOdontogram:
		return "ODONT"
	default:
		return "DOC"
	}
}

func toRomanMonth(m time.Month) string {
	switch m {
	case time.January:
		return "I"
	case time.February:
		return "II"
	case time.March:
		return "III"
	case time.April:
		return "IV"
	case time.May:
		return "V"
	case time.June:
		return "VI"
	case time.July:
		return "VII"
	case time.August:
		return "VIII"
	case time.September:
		return "IX"
	case time.October:
		return "X"
	case time.November:
		return "XI"
	case time.December:
		return "XII"
	default:
		return "I"
	}
}

func formatDocumentNumber(seq int64, prefix, orgCode string, t time.Time) string {
	return fmt.Sprintf("%04d/%s/%s/%s/%d", seq, prefix, orgCode, toRomanMonth(t.Month()), t.Year())
}

func generateVerificationToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
