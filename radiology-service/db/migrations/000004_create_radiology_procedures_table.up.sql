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
