CREATE TABLE allergies (
    id                  CHAR(26) PRIMARY KEY,
    patient_id          CHAR(26)     NOT NULL,
    encounter_id        VARCHAR(26)  NOT NULL DEFAULT '',
    practitioner_id     VARCHAR(64)  NOT NULL DEFAULT '',

    clinical_status     SMALLINT     NOT NULL DEFAULT 1,
    verification_status SMALLINT     NOT NULL DEFAULT 0,
    type                SMALLINT     NOT NULL DEFAULT 0,
    severity            SMALLINT     NOT NULL DEFAULT 0,

    satusehat_id        VARCHAR(64)  NOT NULL DEFAULT '',
    kfa_code            VARCHAR(50)  NOT NULL DEFAULT '',
    snomed_code         VARCHAR(50)  NOT NULL DEFAULT '',
    allergen            VARCHAR(255) NOT NULL,
    reaction            TEXT         NOT NULL DEFAULT '',
    notes               TEXT         NOT NULL DEFAULT '',

    onset_at            TIMESTAMPTZ,
    recorded_at         TIMESTAMPTZ  NOT NULL DEFAULT now(),
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_allergies_patient_id   ON allergies (patient_id);
CREATE INDEX idx_allergies_encounter_id ON allergies (encounter_id) WHERE encounter_id <> '';
