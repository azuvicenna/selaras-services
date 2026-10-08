CREATE TABLE diagnostic_reports (
    id                       CHAR(26) PRIMARY KEY,
    patient_id               CHAR(26)     NOT NULL,
    encounter_id             CHAR(26)     NOT NULL REFERENCES encounters (id),
    service_request_id       VARCHAR(64)  NOT NULL DEFAULT '',
    requester_id             VARCHAR(64)  NOT NULL DEFAULT '',
    practitioner_id          VARCHAR(64)  NOT NULL DEFAULT '',

    status                   SMALLINT     NOT NULL DEFAULT 0,
    category                 SMALLINT     NOT NULL DEFAULT 0,

    satusehat_id             VARCHAR(64)  NOT NULL DEFAULT '',
    report_code              VARCHAR(50)  NOT NULL DEFAULT '',
    report_name              VARCHAR(255) NOT NULL,

    results                  JSONB        NOT NULL DEFAULT '[]'::jsonb,
    conclusion               TEXT         NOT NULL DEFAULT '',
    conclusion_codes         TEXT[]       NOT NULL DEFAULT '{}',

    attachment_urls          TEXT[]       NOT NULL DEFAULT '{}',
    dicom_study_instance_uid VARCHAR(128) NOT NULL DEFAULT '',
    specimen_id              VARCHAR(64)  NOT NULL DEFAULT '',
    amended_from_report_id   VARCHAR(26)  NOT NULL DEFAULT '',

    effective_at             TIMESTAMPTZ,
    issued_at                TIMESTAMPTZ,
    created_at               TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX        idx_diagnostic_reports_patient_id   ON diagnostic_reports (patient_id);
CREATE INDEX        idx_diagnostic_reports_encounter_id ON diagnostic_reports (encounter_id);
CREATE UNIQUE INDEX idx_diagnostic_reports_satusehat_id ON diagnostic_reports (satusehat_id) WHERE satusehat_id <> '';
