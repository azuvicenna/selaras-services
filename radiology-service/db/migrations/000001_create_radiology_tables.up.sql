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

CREATE TABLE radiology_orders (
    id                  CHAR(26)      PRIMARY KEY,
    order_number        VARCHAR(50)   NOT NULL UNIQUE,
    encounter_id        VARCHAR(50)   NOT NULL,
    patient_id          VARCHAR(50)   NOT NULL,
    patient_name        VARCHAR(150)  NOT NULL,
    patient_gender      VARCHAR(20)   NOT NULL DEFAULT '',
    patient_dob         VARCHAR(10)   NOT NULL DEFAULT '',
    doctor_id           VARCHAR(50)   NOT NULL DEFAULT '',
    doctor_name         VARCHAR(150)  NOT NULL DEFAULT '',
    clinical_indication TEXT          NOT NULL DEFAULT '',
    priority            SMALLINT      NOT NULL DEFAULT 1,
    status              SMALLINT      NOT NULL DEFAULT 1,
    total_price         NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_radiology_orders_number ON radiology_orders (order_number);
CREATE INDEX idx_radiology_orders_encounter ON radiology_orders (encounter_id);
CREATE INDEX idx_radiology_orders_patient ON radiology_orders (patient_id);
CREATE INDEX idx_radiology_orders_status ON radiology_orders (status);
CREATE INDEX idx_radiology_orders_created_at ON radiology_orders (created_at);

CREATE TABLE radiology_order_items (
    id         CHAR(26)      PRIMARY KEY,
    order_id   CHAR(26)      NOT NULL REFERENCES radiology_orders(id) ON DELETE CASCADE,
    exam_id    CHAR(26)      NOT NULL REFERENCES radiology_exams(id) ON DELETE RESTRICT,
    exam_code  VARCHAR(50)   NOT NULL,
    exam_name  VARCHAR(150)  NOT NULL,
    modality   SMALLINT      NOT NULL DEFAULT 1,
    price      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    notes      TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_radiology_order_items_order ON radiology_order_items (order_id);
CREATE INDEX idx_radiology_order_items_exam ON radiology_order_items (exam_id);

CREATE TABLE radiology_procedures (
    id                  CHAR(26)     PRIMARY KEY,
    order_id            CHAR(26)     NOT NULL UNIQUE REFERENCES radiology_orders(id) ON DELETE CASCADE,
    room_or_device      VARCHAR(100) NOT NULL DEFAULT '',
    radiographer_id     VARCHAR(50)  NOT NULL DEFAULT '',
    radiographer_name   VARCHAR(150) NOT NULL DEFAULT '',
    contrast_agent_used VARCHAR(100) NOT NULL DEFAULT '',
    series_count        INT          NOT NULL DEFAULT 0,
    instance_count      INT          NOT NULL DEFAULT 0,
    dicom_study_uid     VARCHAR(128) NOT NULL DEFAULT '',
    pacs_viewer_url     TEXT         NOT NULL DEFAULT '',
    notes               TEXT         NOT NULL DEFAULT '',
    scheduled_at        TIMESTAMPTZ,
    performed_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_radiology_procedures_order ON radiology_procedures (order_id);
CREATE INDEX idx_radiology_procedures_dicom_study ON radiology_procedures (dicom_study_uid);

CREATE TABLE radiology_reports (
    id               CHAR(26)     PRIMARY KEY,
    order_id         CHAR(26)     NOT NULL UNIQUE REFERENCES radiology_orders(id) ON DELETE CASCADE,
    order_number     VARCHAR(50)  NOT NULL,
    patient_id       VARCHAR(50)  NOT NULL,
    radiologist_id   VARCHAR(50)  NOT NULL DEFAULT '',
    radiologist_name VARCHAR(150) NOT NULL DEFAULT '',
    clinical_history TEXT         NOT NULL DEFAULT '',
    findings         TEXT         NOT NULL DEFAULT '',
    impression       TEXT         NOT NULL DEFAULT '',
    recommendation   TEXT         NOT NULL DEFAULT '',
    is_critical      BOOLEAN      NOT NULL DEFAULT FALSE,
    is_verified      BOOLEAN      NOT NULL DEFAULT FALSE,
    verified_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_radiology_reports_order ON radiology_reports (order_id);
CREATE INDEX idx_radiology_reports_patient ON radiology_reports (patient_id);
CREATE INDEX idx_radiology_reports_verified ON radiology_reports (is_verified);
CREATE INDEX idx_radiology_reports_critical ON radiology_reports (is_critical);
