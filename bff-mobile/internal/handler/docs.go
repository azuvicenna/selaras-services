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
  <title>Selaras SIMRS Mobile BFF API Docs</title>
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
    "title": "Selaras SIMRS - Mobile Backend For Frontend (BFF) API",
    "description": "API Gateway dan komposit agregasi untuk aplikasi mobile pasien dan tenaga kesehatan (Flutter / React Native) RS Selaras Medika.",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8081",
      "description": "Local Development Server"
    }
  ],
  "paths": {
    "/healthz": {
      "get": {
        "summary": "Health Check",
        "responses": {
          "200": { "description": "Layanan sehat" }
        }
      }
    },
    "/api/v1/auth/login": {
      "post": {
        "summary": "Login Mobile (User/Pasien)",
        "responses": {
          "200": { "description": "Token JWT dan data user" }
        }
      }
    },
    "/api/v1/home": {
      "get": {
        "summary": "Mobile Dashboard Composite",
        "description": "Mengambil agregasi profil pasien, jadwal kunjungan terdekat, status antrean aktif, dan ketersediaan tempat tidur dalam satu panggilan.",
        "responses": {
          "200": { "description": "Dashboard data berhasil diambil" }
        }
      }
    },
    "/api/v1/booking/polyclinics": {
      "get": {
        "summary": "Daftar Poliklinik",
        "responses": {
          "200": { "description": "Daftar poli untuk reservasi" }
        }
      }
    },
    "/api/v1/booking/doctors": {
      "get": {
        "summary": "Daftar Dokter dan Jadwal Praktik",
        "responses": {
          "200": { "description": "Daftar dokter spesialis" }
        }
      }
    },
    "/api/v1/booking/reservations": {
      "post": {
        "summary": "Buat Janji Temu / Reservasi Kunjungan",
        "responses": {
          "201": { "description": "Reservasi berhasil dibuat" }
        }
      }
    },
    "/api/v1/queue/live": {
      "get": {
        "summary": "Status Antrean Live",
        "responses": {
          "200": { "description": "Nomor antrean yang sedang dipanggil dan sisa antrean" }
        }
      }
    },
    "/api/v1/records/encounters": {
      "get": {
        "summary": "Riwayat Kunjungan Medis Pasien",
        "responses": {
          "200": { "description": "Daftar encounter/kunjungan masa lalu" }
        }
      }
    },
    "/api/v1/lab/results": {
      "get": {
        "summary": "Hasil Pemeriksaan Laboratorium Pasien",
        "responses": {
          "200": { "description": "Daftar hasil lab terverifikasi" }
        }
      }
    },
    "/api/v1/radiology/results": {
      "get": {
        "summary": "Hasil Pemeriksaan Radiologi Pasien",
        "responses": {
          "200": { "description": "Daftar hasil rontgen/USG/CT-Scan" }
        }
      }
    },
    "/api/v1/billing/invoices": {
      "get": {
        "summary": "Riwayat Tagihan & Pembayaran",
        "responses": {
          "200": { "description": "Daftar invoice pasien" }
        }
      }
    },
    "/api/v1/beds/availability": {
      "get": {
        "summary": "Informasi Ketersediaan Kamar / Tempat Tidur",
        "responses": {
          "200": { "description": "Kapasitas dan sisa bed per kelas rawat" }
        }
      }
    }
  }
}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, spec)
}
