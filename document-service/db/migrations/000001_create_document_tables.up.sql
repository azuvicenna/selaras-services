CREATE TABLE document_templates (
    id            CHAR(26)     PRIMARY KEY,
    code          VARCHAR(50)  NOT NULL UNIQUE,
    name          VARCHAR(150) NOT NULL,
    type          SMALLINT     NOT NULL DEFAULT 1,
    typst_content TEXT         NOT NULL DEFAULT '',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    version       INT          NOT NULL DEFAULT 1,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_document_templates_code ON document_templates (code);
CREATE INDEX idx_document_templates_type ON document_templates (type);
CREATE INDEX idx_document_templates_is_active ON document_templates (is_active);

CREATE TABLE medical_documents (
    id                 CHAR(26)     PRIMARY KEY,
    document_number    VARCHAR(60)  NOT NULL UNIQUE,
    type               SMALLINT     NOT NULL DEFAULT 1,
    title              VARCHAR(150) NOT NULL,
    encounter_id       VARCHAR(50)  NOT NULL DEFAULT '',
    patient_id         VARCHAR(50)  NOT NULL,
    patient_name       VARCHAR(150) NOT NULL,
    patient_gender     VARCHAR(20)  NOT NULL DEFAULT '',
    patient_dob        VARCHAR(10)  NOT NULL DEFAULT '',
    doctor_id          VARCHAR(50)  NOT NULL DEFAULT '',
    doctor_name        VARCHAR(150) NOT NULL DEFAULT '',
    doctor_sip         VARCHAR(100) NOT NULL DEFAULT '',
    status             SMALLINT     NOT NULL DEFAULT 1,
    payload_json       TEXT         NOT NULL DEFAULT '{}',
    typst_source       TEXT         NOT NULL DEFAULT '',
    verification_token VARCHAR(64)  NOT NULL DEFAULT '',
    qr_code_url        TEXT         NOT NULL DEFAULT '',
    issued_at          TIMESTAMPTZ,
    valid_until        TIMESTAMPTZ,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_medical_documents_number ON medical_documents (document_number);
CREATE INDEX idx_medical_documents_patient ON medical_documents (patient_id);
CREATE INDEX idx_medical_documents_encounter ON medical_documents (encounter_id);
CREATE INDEX idx_medical_documents_type ON medical_documents (type);
CREATE INDEX idx_medical_documents_status ON medical_documents (status);
CREATE INDEX idx_medical_documents_token ON medical_documents (verification_token);
CREATE INDEX idx_medical_documents_created_at ON medical_documents (created_at);

CREATE TABLE document_counters (
    prefix           VARCHAR(20) NOT NULL,
    year             INT         NOT NULL,
    current_sequence BIGINT      NOT NULL DEFAULT 0,
    PRIMARY KEY (prefix, year)
);
