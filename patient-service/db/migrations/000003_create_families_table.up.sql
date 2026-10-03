CREATE TABLE families (
    id                   CHAR(26) PRIMARY KEY,
    patient_id           CHAR(26)     NOT NULL REFERENCES patients (id),
    name                 VARCHAR(255) NOT NULL,
    nik                  VARCHAR(16)  NOT NULL DEFAULT '',
    relation             SMALLINT     NOT NULL DEFAULT 0,
    phone                VARCHAR(20)  NOT NULL DEFAULT '',
    address              TEXT         NOT NULL DEFAULT '',
    is_emergency_contact BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_families_patient_id ON families (patient_id);