package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidChannel      = errors.New("invalid or unsupported notification channel")
	ErrMissingRecipient    = errors.New("recipient address or phone number is required")
	ErrEmptyNotification   = errors.New("notification content or valid template data is required")
	ErrNotificationFailed  = errors.New("failed to send notification")
)

type ChannelType string

const (
	ChannelWhatsApp ChannelType = "WHATSAPP"
	ChannelEmail    ChannelType = "EMAIL"
	ChannelSMS      ChannelType = "SMS"
)

type TemplateType string

const (
	TemplateAppointmentReminder    TemplateType = "APPOINTMENT_REMINDER"
	TemplateRegistrationWelcome    TemplateType = "REGISTRATION_WELCOME"
	TemplateQueueTicketUpdate      TemplateType = "QUEUE_TICKET_UPDATE"
	TemplateLabResultReady         TemplateType = "LAB_RESULT_READY"
	TemplateRadiologyResultReady   TemplateType = "RADIOLOGY_RESULT_READY"
	TemplateBillingInvoice         TemplateType = "BILLING_INVOICE"
	TemplateOTP                    TemplateType = "OTP_CODE"
	TemplateCustom                 TemplateType = "CUSTOM"
)

// NotificationMessage is the event payload transmitted over RabbitMQ.
type NotificationMessage struct {
	ID          string         `json:"id"`
	Channel     ChannelType    `json:"channel"`
	Recipient   string         `json:"recipient"` // Phone (+62...) or Email
	Type        TemplateType   `json:"type"`      // Template type
	Subject     string         `json:"subject,omitempty"` // For Email
	Content     string         `json:"content,omitempty"` // Raw content or fallback
	Data        map[string]any `json:"data,omitempty"`    // Dynamic key-values for template
	ReferenceID string         `json:"reference_id,omitempty"` // PatientID, EncounterID, InvoiceID, etc.
	RetryCount  int            `json:"retry_count,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (m *NotificationMessage) Validate() error {
	m.Recipient = strings.TrimSpace(m.Recipient)
	if m.Recipient == "" {
		return ErrMissingRecipient
	}

	normChannel := ChannelType(strings.ToUpper(strings.TrimSpace(string(m.Channel))))
	switch normChannel {
	case ChannelWhatsApp, ChannelEmail, ChannelSMS:
		m.Channel = normChannel
	default:
		return ErrInvalidChannel
	}

	if strings.TrimSpace(m.Content) == "" && len(m.Data) == 0 && m.Type == "" {
		return ErrEmptyNotification
	}

	return nil
}

// NotificationResult represents the outcome of a notification attempt.
type NotificationResult struct {
	MessageID    string      `json:"message_id"`
	Channel      ChannelType `json:"channel"`
	Recipient    string      `json:"recipient"`
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message,omitempty"`
	Timestamp    time.Time   `json:"timestamp"`
}
