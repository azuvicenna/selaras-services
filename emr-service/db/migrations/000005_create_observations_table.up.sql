CREATE TABLE observations (
    id                          CHAR(26) PRIMARY KEY,
    patient_id                  CHAR(26)     NOT NULL,
    encounter_id                CHAR(26)     NOT NULL REFERENCES encounters (id),
    practitioner_id             VARCHAR(64)  NOT NULL,

    status                      SMALLINT     NOT NULL DEFAULT 0,
    category                    SMALLINT     NOT NULL DEFAULT 0,

    satusehat_id                VARCHAR(64)  NOT NULL DEFAULT '',
    code                        VARCHAR(50)  NOT NULL DEFAULT '',
    name                        VARCHAR(255) NOT NULL,

    value_quantity              DOUBLE PRECISION,
    value_string                TEXT         NOT NULL DEFAULT '',
    value_boolean               BOOLEAN,
    value_code                  VARCHAR(100) NOT NULL DEFAULT '',

    unit                        VARCHAR(50)  NOT NULL DEFAULT '',
    reference_range             VARCHAR(100) NOT NULL DEFAULT '',
    interpretation              VARCHAR(50)  NOT NULL DEFAULT '',

    components                  JSONB        NOT NULL DEFAULT '[]'::jsonb,

    body_site                   VARCHAR(100) NOT NULL DEFAULT '',
    method                      VARCHAR(100) NOT NULL DEFAULT '',
    notes                       TEXT         NOT NULL DEFAULT '',

    amended_from_observation_id VARCHAR(26)  NOT NULL DEFAULT '',

    effective_time              TIMESTAMPTZ  NOT NULL DEFAULT now(),
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_observations_patient_id   ON observations (patient_id);
CREATE INDEX idx_observations_encounter_id ON observations (encounter_id);
CREATE INDEX idx_observations_category     ON observations (patient_id, category);
