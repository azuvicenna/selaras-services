CREATE TABLE encounters (
    id                    CHAR(26) PRIMARY KEY,
    patient_id            CHAR(26)     NOT NULL,
    practitioner_id       VARCHAR(64)  NOT NULL DEFAULT '',
    location_id           VARCHAR(64)  NOT NULL DEFAULT '',

    status                SMALLINT     NOT NULL DEFAULT 1,
    encounter_class       SMALLINT     NOT NULL DEFAULT 0,
    priority              SMALLINT     NOT NULL DEFAULT 0,

    chief_complaint       TEXT         NOT NULL DEFAULT '',

    satusehat_id          VARCHAR(64)  NOT NULL DEFAULT '',
    sep_number            VARCHAR(50)  NOT NULL DEFAULT '',

    parent_encounter_id   VARCHAR(26)  NOT NULL DEFAULT '',
    referral_id           VARCHAR(64)  NOT NULL DEFAULT '',

    service_type          VARCHAR(100) NOT NULL DEFAULT '',
    participant_ids       TEXT[]       NOT NULL DEFAULT '{}',
    discharge_disposition SMALLINT     NOT NULL DEFAULT 0,

    start_time            TIMESTAMPTZ,
    end_time              TIMESTAMPTZ,

    created_at            TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX        idx_encounters_patient_id      ON encounters (patient_id);
CREATE INDEX        idx_encounters_practitioner_id ON encounters (practitioner_id);
CREATE UNIQUE INDEX idx_encounters_satusehat_id    ON encounters (satusehat_id) WHERE satusehat_id <> '';
