CREATE TABLE radiology_exams (
    id                      CHAR(26)      PRIMARY KEY,
    code                    VARCHAR(50)   NOT NULL UNIQUE,
    name                    VARCHAR(150)  NOT NULL,
    modality                SMALLINT      NOT NULL DEFAULT 1,
    body_site               SMALLINT      NOT NULL DEFAULT 1,
    loinc_code              VARCHAR(50)   NOT NULL DEFAULT '',
    requires_contrast       BOOLEAN       NOT NULL DEFAULT FALSE,
    price                   NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    preparation_instructions TEXT         NOT NULL DEFAULT '',
    is_active               BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_radiology_exams_code ON radiology_exams (code);
CREATE INDEX idx_radiology_exams_modality ON radiology_exams (modality);
CREATE INDEX idx_radiology_exams_body_site ON radiology_exams (body_site);
CREATE INDEX idx_radiology_exams_is_active ON radiology_exams (is_active);
CREATE INDEX idx_radiology_exams_loinc_code ON radiology_exams (loinc_code);
