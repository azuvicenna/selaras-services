package processor_test

import (
	"context"
	"strings"
	"testing"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/processor"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/provider"
	"github.com/azuvicenna/selaras-services/notification-worker/internal/template"
)

func TestProcessor_Process_WhatsAppSuccess(t *testing.T) {
	waMock := provider.NewMockWhatsAppProvider()
	emailMock := provider.NewMockEmailProvider()
	renderer := template.NewRenderer()

	proc := processor.NewProcessor(renderer, waMock, emailMock, nil)

	msg := domain.NotificationMessage{
		ID:        "msg-01",
		Channel:   domain.ChannelWhatsApp,
		Recipient: "+628123456789",
		Type:      domain.TemplateOTP,
		Data: map[string]any{
			"otp": "456789",
		},
	}

	res, err := proc.Process(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}

	messages := waMock.GetMessages()
	if len(messages) != 1 {
		t.Fatalf("expected 1 sent message, got %d", len(messages))
	}
	if !strings.Contains(messages[0].Message, "456789") {
		t.Errorf("expected message to contain OTP 456789, got: %s", messages[0].Message)
	}
}

func TestProcessor_Process_EmailSuccess(t *testing.T) {
	waMock := provider.NewMockWhatsAppProvider()
	emailMock := provider.NewMockEmailProvider()
	renderer := template.NewRenderer()

	proc := processor.NewProcessor(renderer, waMock, emailMock, nil)

	msg := domain.NotificationMessage{
		ID:        "msg-02",
		Channel:   domain.ChannelEmail,
		Recipient: "pasien@example.com",
		Type:      domain.TemplateBillingInvoice,
		Data: map[string]any{
			"patient_name":   "Citra",
			"invoice_number": "INV-00123",
			"amount":         "500.000",
		},
	}

	res, err := proc.Process(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success")
	}

	emails := emailMock.GetEmails()
	if len(emails) != 1 {
		t.Fatalf("expected 1 sent email, got %d", len(emails))
	}
	if !strings.Contains(emails[0].Subject, "INV-00123") {
		t.Errorf("expected subject to contain invoice number, got: %s", emails[0].Subject)
	}
	if !strings.Contains(emails[0].Body, "500.000") {
		t.Errorf("expected body to contain amount, got: %s", emails[0].Body)
	}
}

func TestProcessor_Process_EmptyRecipient(t *testing.T) {
	proc := processor.NewProcessor(nil, provider.NewMockWhatsAppProvider(), nil, nil)

	msg := domain.NotificationMessage{
		ID:        "msg-03",
		Channel:   domain.ChannelWhatsApp,
		Recipient: "", // Invalid
	}

	_, err := proc.Process(context.Background(), msg)
	if err == nil {
		t.Fatal("expected error on empty recipient, got nil")
	}
}
