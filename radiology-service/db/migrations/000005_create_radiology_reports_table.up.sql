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
