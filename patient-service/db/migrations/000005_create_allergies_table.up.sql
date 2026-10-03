CREATE TABLE allergies (
    id         CHAR(26) PRIMARY KEY,
    patient_id CHAR(26)     NOT NULL REFERENCES patients (id),
    type       SMALLINT     NOT NULL DEFAULT 0,
    allergen   VARCHAR(255) NOT NULL,
    severity   SMALLINT     NOT NULL DEFAULT 0,
    reaction   TEXT         NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_allergies_patient_id ON allergies (patient_id);