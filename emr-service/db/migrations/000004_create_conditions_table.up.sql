CREATE TABLE conditions (
    id                  CHAR(26) PRIMARY KEY,
    patient_id          CHAR(26)     NOT NULL,
    encounter_id        CHAR(26)     NOT NULL REFERENCES encounters (id),
    practitioner_id     VARCHAR(64)  NOT NULL,
    clinical_note_id    VARCHAR(26)  NOT NULL DEFAULT '',

    clinical_status     SMALLINT     NOT NULL DEFAULT 1,
    verification_status SMALLINT     NOT NULL DEFAULT 0,
    category            SMALLINT     NOT NULL DEFAULT 0,

    is_primary          BOOLEAN      NOT NULL DEFAULT FALSE,
    satusehat_id        VARCHAR(64)  NOT NULL DEFAULT '',
    icd10_code          VARCHAR(20)  NOT NULL DEFAULT '',
    snomed_code         VARCHAR(50)  NOT NULL DEFAULT '',
    name                VARCHAR(255) NOT NULL,
    severity            VARCHAR(50)  NOT NULL DEFAULT '',
    notes               TEXT         NOT NULL DEFAULT '',

    onset_at            TIMESTAMPTZ,
    abatement_at        TIMESTAMPTZ,
    recorded_at         TIMESTAMPTZ  NOT NULL DEFAULT now(),
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX        idx_conditions_patient_id   ON conditions (patient_id);
CREATE INDEX        idx_conditions_encounter_id ON conditions (encounter_id);
CREATE INDEX        idx_conditions_icd10_code   ON conditions (icd10_code) WHERE icd10_code <> '';
CREATE UNIQUE INDEX idx_conditions_satusehat_id ON conditions (satusehat_id) WHERE satusehat_id <> '';
