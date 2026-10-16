package provider

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/config"
)

type EmailProvider interface {
	Send(ctx context.Context, to, subject, htmlBody string) error
}

// 1. Mock Email Provider for testing & development
type SentEmail struct {
	To        string
	Subject   string
	Body      string
	Timestamp time.Time
}

type MockEmailProvider struct {
	mu     sync.RWMutex
	emails []SentEmail
}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{
		emails: make([]SentEmail, 0),
	}
}

func (m *MockEmailProvider) Send(ctx context.Context, to, subject, htmlBody string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emails = append(m.emails, SentEmail{
		To:        to,
		Subject:   subject,
		Body:      htmlBody,
		Timestamp: time.Now(),
	})
	return nil
}

func (m *MockEmailProvider) GetEmails() []SentEmail {
	m.mu.RLock()
	defer m.mu.RUnlock()
	copied := make([]SentEmail, len(m.emails))
	copy(copied, m.emails)
	return copied
}

// 2. SMTP Email Provider
type smtpEmailProvider struct {
	host     string
	port     int
	user     string
	password string
	from     string
	fromName string
}

func NewSMTPEmailProvider(cfg *config.Config) EmailProvider {
	return &smtpEmailProvider{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		user:     cfg.SMTPUser,
		password: cfg.SMTPPassword,
		from:     cfg.SMTPFrom,
		fromName: cfg.SMTPFromName,
	}
}

func (p *smtpEmailProvider) Send(ctx context.Context, to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)

	fromHeader := p.from
	if p.fromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", p.fromName, p.from)
	}

	headers := make(map[string]string)
	headers["From"] = fromHeader
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(htmlBody)

	var auth smtp.Auth
	if p.user != "" && p.password != "" {
		auth = smtp.PlainAuth("", p.user, p.password, p.host)
	}

	return smtp.SendMail(addr, auth, p.from, []string{to}, []byte(msgBuilder.String()))
}
