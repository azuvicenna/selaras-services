package handler

import (
	"fmt"
	"net/http"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Selaras SIMRS B2B Partner Gateway API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5/favicon-32x32.png" />
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin:0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, html)
}

func (h *DocsHandler) OpenAPISpec(w http.ResponseWriter, r *http.Request) {
	spec := `{
  "openapi": "3.0.3",
  "info": {
    "title": "Selaras SIMRS - B2B Partner Gateway API",
    "description": "API Gateway Enterprise untuk mitra asuransi swasta (Prudential, Allianz, AdMedika, dll.) dan korporat MCU RS Selaras Medika.",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8082",
      "description": "Local B2B Gateway Server"
    }
  ],
  "components": {
    "securitySchemes": {
      "PartnerAuth": {
        "type": "apiKey",
        "in": "header",
        "name": "X-API-Key",
        "description": "API Key B2B Partner"
      },
      "PartnerID": {
        "type": "apiKey",
        "in": "header",
        "name": "X-Partner-Id",
        "description": "Partner Identifier Code (e.g. PRU-INDONESIA)"
      }
    }
  },
  "paths": {
    "/healthz": {
      "get": {
        "summary": "Health Check",
        "responses": {
          "200": { "description": "Gateway sehat" }
        }
      }
    },
    "/api/v1/b2b/insurance/eligibility": {
      "post": {
        "summary": "Verifikasi Eligibilitas Peserta Asuransi",
        "security": [{ "PartnerAuth": [], "PartnerID": [] }],
        "responses": {
          "200": { "description": "Status eligibilitas dan limit plafon" },
          "401": { "description": "Unauthorized" }
        }
      }
    },
    "/api/v1/b2b/insurance/pre-authorization": {
      "post": {
        "summary": "Pengajuan Pre-Authorization Tindakan / Rawat Inap",
        "security": [{ "PartnerAuth": [], "PartnerID": [] }],
        "responses": {
          "201": { "description": "Pre-authorization disetujui" }
        }
      }
    },
    "/api/v1/b2b/insurance/claims": {
      "post": {
        "summary": "Pengajuan Klaim Tagihan Pasien",
        "security": [{ "PartnerAuth": [], "PartnerID": [] }],
        "responses": {
          "201": { "description": "Klaim berhasil disubmit dan diverifikasi" }
        }
      },
      "get": {
        "summary": "Pengecekan Status Klaim",
        "security": [{ "PartnerAuth": [], "PartnerID": [] }],
        "responses": {
          "200": { "description": "Detail status proses klaim" }
        }
      }
    },
    "/api/v1/b2b/corporate/mcu/register": {
      "post": {
        "summary": "Pendaftaran Batch MCU Karyawan Korporat",
        "security": [{ "PartnerAuth": [], "PartnerID": [] }],
        "responses": {
          "201": { "description": "Batch pendaftaran MCU berhasil" }
        }
      }
    },
    "/api/v1/b2b/corporate/mcu/batches": {
      "get": {
        "summary": "Monitoring Status Batch MCU",
        "security": [{ "PartnerAuth": [], "PartnerID": [] }],
        "responses": {
          "200": { "description": "Status progres batch MCU" }
        }
      }
    },
    "/api/v1/b2b/webhooks": {
      "post": {
        "summary": "Pendaftaran Webhook Event Callback Mitra",
        "security": [{ "PartnerAuth": [], "PartnerID": [] }],
        "responses": {
          "201": { "description": "Webhook berhasil didaftarkan" }
        }
      },
      "get": {
        "summary": "Daftar Webhook Aktif Mitra",
        "security": [{ "PartnerAuth": [], "PartnerID": [] }],
        "responses": {
          "200": { "description": "Daftar webhook milik mitra" }
        }
      }
    }
  }
}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, spec)
}
