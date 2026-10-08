CREATE TABLE clinical_notes (
    id                    CHAR(26) PRIMARY KEY,
    patient_id            CHAR(26)    NOT NULL,
    encounter_id          CHAR(26)    NOT NULL REFERENCES encounters (id),
    practitioner_id       VARCHAR(64) NOT NULL,

    status                SMALLINT    NOT NULL DEFAULT 1,
    note_type             SMALLINT    NOT NULL DEFAULT 0,

    subjective            TEXT        NOT NULL DEFAULT '',
    objective             TEXT        NOT NULL DEFAULT '',
    assessment            TEXT        NOT NULL DEFAULT '',
    plan                  TEXT        NOT NULL DEFAULT '',
    free_text_note        TEXT        NOT NULL DEFAULT '',

    amendment_reason      TEXT        NOT NULL DEFAULT '',
    amended_from_note_id  VARCHAR(26) NOT NULL DEFAULT '',
    cosigner_id           VARCHAR(64) NOT NULL DEFAULT '',
    cosigned_at           TIMESTAMPTZ,

    practitioner_role     VARCHAR(50) NOT NULL DEFAULT '',
    unit_id               VARCHAR(64) NOT NULL DEFAULT '',
    is_confidential       BOOLEAN     NOT NULL DEFAULT FALSE,

    primary_icd10_code    VARCHAR(20) NOT NULL DEFAULT '',
    secondary_icd10_codes TEXT[]      NOT NULL DEFAULT '{}',
    icd9cm_codes          TEXT[]      NOT NULL DEFAULT '{}',

    attachment_urls       TEXT[]      NOT NULL DEFAULT '{}',
    observation_ids       TEXT[]      NOT NULL DEFAULT '{}',

    satusehat_id          VARCHAR(64) NOT NULL DEFAULT '',
    digital_signature     TEXT        NOT NULL DEFAULT '',
    signer_id             VARCHAR(64) NOT NULL DEFAULT '',
    signed_at             TIMESTAMPTZ,

    recorded_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX        idx_clinical_notes_patient_id   ON clinical_notes (patient_id);
CREATE INDEX        idx_clinical_notes_encounter_id ON clinical_notes (encounter_id);
CREATE UNIQUE INDEX idx_clinical_notes_satusehat_id ON clinical_notes (satusehat_id) WHERE satusehat_id <> '';
