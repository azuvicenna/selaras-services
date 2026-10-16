package template

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/azuvicenna/selaras-services/notification-worker/internal/domain"
)

type Renderer interface {
	Render(msg *domain.NotificationMessage) (subject string, body string, err error)
}

type templateRenderer struct{}

func NewRenderer() Renderer {
	return &templateRenderer{}
}

func (r *templateRenderer) Render(msg *domain.NotificationMessage) (string, string, error) {
	data := normalizeMapKeys(msg.Data)

	subject := msg.Subject
	var textBody string

	switch msg.Type {
	case domain.TemplateAppointmentReminder:
		if subject == "" {
			subject = "Pengingat Janji Temu Dokter - RS Selaras Medika"
		}
		textBody = fmt.Sprintf(
			"Halo %s,\nIni adalah pengingat jadwal konsultasi Anda di *RS Selaras Medika*:\n🩺 Dokter: %s\n🏥 Poli: %s\n📅 Waktu: %s\n🔢 No. Antrean: %s\n\nMohon hadir 15 menit sebelum waktu pemeriksaan. Terima kasih.",
			getString(data, "patient_name", "Pasien"),
			getString(data, "practitioner_name", "Dokter"),
			getString(data, "poli_name", "Poli"),
			getString(data, "schedule_time", "-"),
			getString(data, "queue_number", "-"),
		)

	case domain.TemplateRegistrationWelcome:
		noRM := getString(data, "medical_record_number", "-")
		if subject == "" {
			subject = fmt.Sprintf("Selamat Datang di RS Selaras Medika - No. RM %s", noRM)
		}
		textBody = fmt.Sprintf(
			"Halo %s,\nPendaftaran pasien baru Anda di *RS Selaras Medika* telah berhasil.\nNomor Rekam Medis (No. RM) Anda: *%s*.\nGunakan nomor ini untuk reservasi dan riwayat kunjungan berikutnya.",
			getString(data, "patient_name", "Pasien"),
			noRM,
		)

	case domain.TemplateQueueTicketUpdate:
		qNum := getString(data, "queue_number", "-")
		if subject == "" {
			subject = fmt.Sprintf("Panggilan Antrean - %s", qNum)
		}
		textBody = fmt.Sprintf(
			"Panggilan Antrean *%s*:\nPasien atas nama *%s* dipersilakan menuju ke *%s*.",
			qNum,
			getString(data, "patient_name", "Pasien"),
			getString(data, "destination", "Ruang Pemeriksaan"),
		)

	case domain.TemplateLabResultReady:
		labNo := getString(data, "lab_number", "-")
		if subject == "" {
			subject = fmt.Sprintf("Hasil Pemeriksaan Laboratorium Selesai - %s", labNo)
		}
		textBody = fmt.Sprintf(
			"Halo %s,\nHasil pemeriksaan laboratorium Anda (No. Lab: *%s*) telah selesai divalidasi.\nUnduh berkas hasil: %s",
			getString(data, "patient_name", "Pasien"),
			labNo,
			getString(data, "download_url", "https://portal.selaras.id"),
		)

	case domain.TemplateRadiologyResultReady:
		modality := getString(data, "modality", "Radiologi")
		if subject == "" {
			subject = fmt.Sprintf("Hasil Pemeriksaan Radiologi Selesai - %s", modality)
		}
		textBody = fmt.Sprintf(
			"Halo %s,\nHasil pemeriksaan radiologi (%s) telah selesai diekspertise oleh Dokter Spesialis Radiologi.\nUnduh hasil: %s",
			getString(data, "patient_name", "Pasien"),
			modality,
			getString(data, "download_url", "https://portal.selaras.id"),
		)

	case domain.TemplateBillingInvoice:
		invNo := getString(data, "invoice_number", "-")
		if subject == "" {
			subject = fmt.Sprintf("Bukti Pembayaran Tagihan RS Selaras Medika - %s", invNo)
		}
		textBody = fmt.Sprintf(
			"Halo %s,\nPembayaran tagihan dengan No. Invoice *%s* sejumlah *Rp %s* telah berhasil kami terima dan verifikasi.\nTerima kasih.",
			getString(data, "patient_name", "Pasien"),
			invNo,
			getString(data, "amount", "0"),
		)

	case domain.TemplateOTP:
		otp := getString(data, "otp", "000000")
		if subject == "" {
			subject = "Kode OTP Verifikasi - RS Selaras Medika"
		}
		textBody = fmt.Sprintf(
			"Kode OTP verifikasi Selaras Anda adalah: *%s*.\nBerlaku selama 5 menit. JANGAN berikan kode ini kepada siapa pun.",
			otp,
		)

	default: // TemplateCustom or raw
		if subject == "" {
			subject = "Pemberitahuan dari RS Selaras Medika"
		}
		textBody = msg.Content
		if textBody == "" {
			textBody = getString(data, "message", "Pemberitahuan dari RS Selaras Medika")
		}
	}

	if msg.Channel == domain.ChannelEmail {
		htmlBody := wrapHTMLTemplate(subject, textBody)
		return subject, htmlBody, nil
	}

	return subject, textBody, nil
}

func wrapHTMLTemplate(subject, textContent string) string {
	paragraphs := strings.Split(textContent, "\n")
	var htmlParagraphs strings.Builder
	for _, p := range paragraphs {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			// Convert *bold* into <strong>bold</strong>
			formatted := template.HTMLEscapeString(trimmed)
			formatted = formatBold(formatted)
			htmlParagraphs.WriteString(fmt.Sprintf("<p style='margin: 8px 0; color: #374151;'>%s</p>", formatted))
		}
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{.Subject}}</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f3f4f6; margin: 0; padding: 20px;">
  <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);">
    <div style="background-color: #0d9488; padding: 20px 24px; color: #ffffff;">
      <h2 style="margin: 0; font-size: 20px; font-weight: 600;">RS Selaras Medika</h2>
    </div>
    <div style="padding: 24px;">
      <h3 style="margin-top: 0; color: #111827; font-size: 18px;">{{.Subject}}</h3>
      {{.Content}}
    </div>
    <div style="background-color: #f9fafb; padding: 16px 24px; border-top: 1px solid #e5e7eb; font-size: 12px; color: #6b7280; text-align: center;">
      Pesan ini dikirim secara otomatis oleh Sistem Rumah Sakit Selaras. Mohon tidak membalas email ini.
    </div>
  </div>
</body>
</html>`

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return textContent
	}

	var buf bytes.Buffer
	_ = t.Execute(&buf, map[string]any{
		"Subject": subject,
		"Content": template.HTML(htmlParagraphs.String()),
	})

	return buf.String()
}

func formatBold(s string) string {
	parts := strings.Split(s, "*")
	if len(parts) < 3 {
		return s
	}
	var res strings.Builder
	for i, part := range parts {
		if i%2 == 1 {
			res.WriteString("<strong>" + part + "</strong>")
		} else {
			res.WriteString(part)
		}
	}
	return res.String()
}

func normalizeMapKeys(m map[string]any) map[string]any {
	res := make(map[string]any, len(m))
	for k, v := range m {
		normKey := strings.ToLower(strings.TrimSpace(k))
		normKey = strings.ReplaceAll(normKey, "-", "_")
		res[normKey] = v
	}
	return res
}

func getString(m map[string]any, key, defaultVal string) string {
	if val, ok := m[key]; ok && val != nil {
		return fmt.Sprintf("%v", val)
	}
	return defaultVal
}
