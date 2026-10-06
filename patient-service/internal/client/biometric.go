package client

import (
	"crypto/subtle"
	"errors"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

var _ domain.BiometricMatcher = (*DefaultBiometricMatcher)(nil)

// DefaultBiometricMatcher mengimplementasikan interface pencocokan biometrik.
// Menggunakan constant-time comparison untuk keamanan dan menyediakan titik ekstensi untuk SDK vendor (Neurotechnology, SecuGen, dll).
type DefaultBiometricMatcher struct {
	threshold float64
}

func NewDefaultBiometricMatcher(threshold float64) *DefaultBiometricMatcher {
	if threshold <= 0 {
		threshold = 0.75
	}
	return &DefaultBiometricMatcher{threshold: threshold}
}

// Match membandingkan template biometrik terdaftar dengan sampel sidik jari yang dikirimkan.
func (m *DefaultBiometricMatcher) Match(template, sample []byte) (bool, error) {
	if len(template) == 0 || len(sample) == 0 {
		return false, errors.New("template and sample data cannot be empty")
	}

	// Constant-time byte comparison untuk mencegah serangan timing attack
	if subtle.ConstantTimeCompare(template, sample) == 1 {
		return true, nil
	}

	// Catatan integrasi: Jika menggunakan engine pencocokan berbasis minutiae/feature score (e.g. SDK BioMini / Griaule),
	// hitung similarity score di sini dan bandingkan dengan m.threshold.
	return false, nil
}
