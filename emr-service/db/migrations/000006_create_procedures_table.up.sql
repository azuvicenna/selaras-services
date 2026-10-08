CREATE TABLE procedures (
    id                  CHAR(26) PRIMARY KEY,
    patient_id          CHAR(26)     NOT NULL,
    encounter_id        CHAR(26)     NOT NULL REFERENCES encounters (id),
    practitioner_id     VARCHAR(64)  NOT NULL,
    reason_condition_id VARCHAR(26)  NOT NULL DEFAULT '',
    clinical_note_id    VARCHAR(26)  NOT NULL DEFAULT '',

    status              SMALLINT     NOT NULL DEFAULT 0,
    category            SMALLINT     NOT NULL DEFAULT 0,

    satusehat_id        VARCHAR(64)  NOT NULL DEFAULT '',
    snomed_code         VARCHAR(50)  NOT NULL DEFAULT '',
    icd9cm_code         VARCHAR(20)  NOT NULL DEFAULT '',
    procedure_name      VARCHAR(255) NOT NULL,

    performers          JSONB        NOT NULL DEFAULT '[]'::jsonb,

    body_site           VARCHAR(100) NOT NULL DEFAULT '',
    outcome             TEXT         NOT NULL DEFAULT '',
    complications       TEXT[]       NOT NULL DEFAULT '{}',
    focal_device_ids    TEXT[]       NOT NULL DEFAULT '{}',
    notes               TEXT         NOT NULL DEFAULT '',

    performed_start     TIMESTAMPTZ,
    performed_end       TIMESTAMPTZ,

    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_procedures_patient_id   ON procedures (patient_id);
CREATE INDEX idx_procedures_encounter_id ON procedures (encounter_id);
