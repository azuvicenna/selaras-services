CREATE TABLE documents (
    id              CHAR(26) PRIMARY KEY,
    patient_id      CHAR(26)     NOT NULL REFERENCES patients (id),
    type            SMALLINT     NOT NULL DEFAULT 0,
    document_number VARCHAR(100) NOT NULL DEFAULT '',
    file_url        TEXT         NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_documents_patient_id ON documents (patient_id);