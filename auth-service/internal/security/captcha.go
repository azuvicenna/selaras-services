package security

import (
	"time"

	"github.com/mojocn/base64Captcha"

	"github.com/azuvicenna/selaras-services/auth-service/internal/domain"
)

var _ domain.CaptchaService = (*Base64CaptchaService)(nil)

type Base64CaptchaService struct {
	captcha *base64Captcha.Captcha
}

func NewBase64CaptchaService(expiry time.Duration) *Base64CaptchaService {
	// Standard 5-digit image captcha with noise dots
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	store := base64Captcha.NewMemoryStore(10240, expiry)
	cp := base64Captcha.NewCaptcha(driver, store)
	return &Base64CaptchaService{captcha: cp}
}

func (s *Base64CaptchaService) Generate() (id string, imageBase64 string, err error) {
	id, b64s, _, err := s.captcha.Generate()
	if err != nil {
		return "", "", err
	}
	return id, b64s, nil
}

func (s *Base64CaptchaService) Verify(id, code string) bool {
	if id == "" || code == "" {
		return false
	}
	// Verify and clear (single-use challenge)
	return s.captcha.Verify(id, code, true)
}
