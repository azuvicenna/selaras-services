package template_test

import (
	"strings"
	"testing"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/template"
)

func TestRenderer_RenderAppointmentReminder_WhatsApp(t *testing.T) {
	r := template.NewRenderer()
	msg := &domain.NotificationMessage{
		Channel: domain.ChannelWhatsApp,
		Type:    domain.TemplateAppointmentReminder,
		Data: map[string]any{
			"patient_name":      "Budi Santoso",
			"practitioner_name": "dr. Andi, Sp.PD",
			"poli_name":         "Poli Penyakit Dalam",
			"schedule_time":     "Senin, 19 Okt 2026, 09:00 WIB",
			"queue_number":      "A-012",
		},
	}

	subj, body, err := r.Render(msg)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	if subj == "" {
		t.Error("expected non-empty subject")
	}
	if !strings.Contains(body, "Budi Santoso") {
		t.Errorf("expected body to contain patient name, got: %s", body)
	}
	if !strings.Contains(body, "dr. Andi, Sp.PD") {
		t.Errorf("expected body to contain practitioner name, got: %s", body)
	}
	if !strings.Contains(body, "A-012") {
		t.Errorf("expected body to contain queue number, got: %s", body)
	}
}

func TestRenderer_RenderEmail_HTML(t *testing.T) {
	r := template.NewRenderer()
	msg := &domain.NotificationMessage{
		Channel: domain.ChannelEmail,
		Type:    domain.TemplateRegistrationWelcome,
		Data: map[string]any{
			"patient_name":          "Citra Dewi",
			"medical_record_number": "RM-2026-9999",
		},
	}

	subj, body, err := r.Render(msg)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	if !strings.Contains(subj, "RM-2026-9999") {
		t.Errorf("expected subject to contain RM number, got: %s", subj)
	}
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Errorf("expected HTML email body, got: %s", body)
	}
	if !strings.Contains(body, "<strong>RM-2026-9999</strong>") {
		t.Errorf("expected bolded RM number in HTML, got: %s", body)
	}
}

func TestRenderer_RenderOTP(t *testing.T) {
	r := template.NewRenderer()
	msg := &domain.NotificationMessage{
		Channel: domain.ChannelWhatsApp,
		Type:    domain.TemplateOTP,
		Data: map[string]any{
			"otp": "839201",
		},
	}

	_, body, err := r.Render(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(body, "839201") {
		t.Errorf("expected body to contain OTP, got: %s", body)
	}
}
