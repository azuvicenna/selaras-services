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
  <title>Selaras SIMRS Web BFF API Docs</title>
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
    "title": "Selaras SIMRS - Web Backend For Frontend (BFF) API",
    "description": "API Gateway dan komposit agregasi untuk aplikasi Web SIMRS Selaras (Front Office, Dokter, Perawat, Farmasi, Laboratorium, Radiologi, Kasir, dan Manajemen).",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "Local Development Server"
    }
  ],
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  },
  "security": [
    {
      "BearerAuth": []
    }
  ],
  "paths": {
    "/healthz": {
      "get": {
        "summary": "Health Check Endpoint",
        "security": [],
        "responses": { "200": { "description": "OK" } }
      }
    },
    "/api/v1/auth/captcha": {
      "get": {
        "summary": "Generate Captcha",
        "security": [],
        "responses": { "200": { "description": "Captcha image and ID" } }
      }
    },
    "/api/v1/auth/login": {
      "post": {
        "summary": "Login User & Get JWT",
        "security": [],
        "responses": { "200": { "description": "JWT Access & Refresh Token" } }
      }
    },
    "/api/v1/auth/me": {
      "get": {
        "summary": "Get Current Authenticated User Profile",
        "responses": { "200": { "description": "User details" } }
      }
    },
    "/api/v1/patients": {
      "get": {
        "summary": "List / Search Patients",
        "responses": { "200": { "description": "List of patients" } }
      },
      "post": {
        "summary": "Register New Patient",
        "responses": { "201": { "description": "Patient created" } }
      }
    },
    "/api/v1/patients/{id}": {
      "get": {
        "summary": "Get Patient by ID",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Patient details" } }
      }
    },
    "/api/v1/practitioners": {
      "get": {
        "summary": "List Doctors & Healthcare Practitioners",
        "responses": { "200": { "description": "List of practitioners" } }
      }
    },
    "/api/v1/appointments": {
      "get": {
        "summary": "List Appointments / Booking",
        "responses": { "200": { "description": "List of appointments" } }
      },
      "post": {
        "summary": "Create Appointment",
        "responses": { "201": { "description": "Appointment created" } }
      }
    },
    "/api/v1/admissions": {
      "get": {
        "summary": "List Patient Admissions",
        "responses": { "200": { "description": "List of admissions" } }
      },
      "post": {
        "summary": "Create Admission (Pendaftaran Poli/Rawat Inap/IGD)",
        "responses": { "201": { "description": "Admission created" } }
      }
    },
    "/api/v1/emr/encounters": {
      "post": {
        "summary": "Start Clinical Encounter",
        "responses": { "201": { "description": "Encounter started" } }
      }
    },
    "/api/v1/emr/clinical-notes": {
      "post": {
        "summary": "Create SOAP / CPPT Clinical Note",
        "responses": { "201": { "description": "Clinical note created" } }
      }
    },
    "/api/v1/billing/invoices": {
      "get": {
        "summary": "List Invoices",
        "responses": { "200": { "description": "List of invoices" } }
      },
      "post": {
        "summary": "Generate Invoice",
        "responses": { "201": { "description": "Invoice created" } }
      }
    },
    "/api/v1/pharmacy/medications": {
      "get": {
        "summary": "List Medication Catalog",
        "responses": { "200": { "description": "List of medications" } }
      }
    },
    "/api/v1/laboratory/orders": {
      "get": {
        "summary": "List Laboratory Orders",
        "responses": { "200": { "description": "List of lab orders" } }
      },
      "post": {
        "summary": "Create Laboratory Order",
        "responses": { "201": { "description": "Lab order created" } }
      }
    },
    "/api/v1/radiology/orders": {
      "get": {
        "summary": "List Radiology Orders",
        "responses": { "200": { "description": "List of radiology orders" } }
      },
      "post": {
        "summary": "Create Radiology Order",
        "responses": { "201": { "description": "Radiology order created" } }
      }
    },
    "/api/v1/documents/templates": {
      "get": {
        "summary": "List 20 Typst Medical Document Templates",
        "responses": { "200": { "description": "List of templates" } }
      }
    },
    "/api/v1/documents/{id}/download": {
      "get": {
        "summary": "Download Typst Generated PDF Document",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Binary PDF stream" } }
      }
    },
    "/api/v1/master/departments": {
      "get": {
        "summary": "List Hospital Departments / Clinics",
        "responses": { "200": { "description": "List of departments" } }
      }
    },
    "/api/v1/master/hospital-profile": {
      "get": {
        "summary": "Get Hospital Profile",
        "responses": { "200": { "description": "Hospital profile details" } }
      }
    },
    "/api/v1/master/icd10": {
      "get": {
        "summary": "Search & List ICD-10 Diagnoses",
        "responses": { "200": { "description": "List of ICD-10 codes" } }
      }
    },
    "/api/v1/master/icd9cm": {
      "get": {
        "summary": "Search & List ICD-9-CM Procedures",
        "responses": { "200": { "description": "List of ICD-9-CM codes" } }
      }
    },
    "/api/v1/master/provinces": {
      "get": {
        "summary": "List Indonesian Provinces",
        "responses": { "200": { "description": "List of provinces" } }
      }
    },
    "/api/v1/master/cities": {
      "get": {
        "summary": "List Cities by Province Code",
        "responses": { "200": { "description": "List of cities" } }
      }
    },
    "/api/v1/master/districts": {
      "get": {
        "summary": "List Districts by City Code",
        "responses": { "200": { "description": "List of districts" } }
      }
    },
    "/api/v1/master/villages": {
      "get": {
        "summary": "List Villages / Kelurahan by District Code",
        "responses": { "200": { "description": "List of villages" } }
      }
    },
    "/api/v1/composite/encounters/{id}/summary": {
      "get": {
        "summary": "Composite Encounter Summary (Encounter + Patient + Doctor + SOAP + Vitals + Labs + Rad + Rx)",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Aggregated clinical dashboard summary" } }
      }
    },
    "/api/v1/composite/patients/{id}/overview": {
      "get": {
        "summary": "Composite Patient 360 Overview (Demographics + Allergies + Visits + Pending Invoices)",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Aggregated patient overview" } }
      }
    }
  }
}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, spec)
}
