CREATE TABLE medication_requests (
    id                        CHAR(26) PRIMARY KEY,
    patient_id                CHAR(26)         NOT NULL,
    encounter_id              CHAR(26)         NOT NULL REFERENCES encounters (id),
    practitioner_id           VARCHAR(64)      NOT NULL,
    condition_id              VARCHAR(26)      NOT NULL DEFAULT '',

    status                    SMALLINT         NOT NULL DEFAULT 1,
    priority                  SMALLINT         NOT NULL DEFAULT 0,
    category                  SMALLINT         NOT NULL DEFAULT 0,

    satusehat_id              VARCHAR(64)      NOT NULL DEFAULT '',
    kfa_code                  VARCHAR(50)      NOT NULL DEFAULT '',
    medication_code           VARCHAR(50)      NOT NULL DEFAULT '',
    medication_name           VARCHAR(255)     NOT NULL DEFAULT '',

    is_compound               BOOLEAN          NOT NULL DEFAULT FALSE,
    compound_name             VARCHAR(255)     NOT NULL DEFAULT '',
    ingredients               JSONB            NOT NULL DEFAULT '[]'::jsonb,

    dosage_instruction        VARCHAR(255)     NOT NULL DEFAULT '',
    route                     VARCHAR(50)      NOT NULL DEFAULT '',
    patient_instruction       TEXT             NOT NULL DEFAULT '',
    duration_in_days          INTEGER          NOT NULL DEFAULT 0,

    dispense_quantity         DOUBLE PRECISION NOT NULL DEFAULT 0,
    dispense_unit             VARCHAR(50)      NOT NULL DEFAULT '',
    number_of_refills_allowed INTEGER          NOT NULL DEFAULT 0,
    substitution_allowed      BOOLEAN          NOT NULL DEFAULT TRUE,

    reason_code               VARCHAR(20)      NOT NULL DEFAULT '',

    authored_on               TIMESTAMPTZ      NOT NULL DEFAULT now(),
    created_at                TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at                TIMESTAMPTZ      NOT NULL DEFAULT now()
);

CREATE INDEX idx_medication_requests_patient_id   ON medication_requests (patient_id);
CREATE INDEX idx_medication_requests_encounter_id ON medication_requests (encounter_id);
