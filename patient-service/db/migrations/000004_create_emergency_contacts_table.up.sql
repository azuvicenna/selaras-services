CREATE TABLE emergency_contacts (
    id           CHAR(26) PRIMARY KEY,
    patient_id   CHAR(26)     NOT NULL REFERENCES patients (id),
    name         VARCHAR(255) NOT NULL,
    relationship VARCHAR(100) NOT NULL DEFAULT '',
    phone        VARCHAR(20)  NOT NULL,
    address      TEXT         NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_emergency_contacts_patient_id ON emergency_contacts (patient_id);