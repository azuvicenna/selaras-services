CREATE TABLE lab_results (
    id               CHAR(26)     PRIMARY KEY,
    order_id         CHAR(26)     NOT NULL UNIQUE REFERENCES lab_orders(id) ON DELETE CASCADE,
    order_number     VARCHAR(50)  NOT NULL,
    patient_id       VARCHAR(50)  NOT NULL,
    analyst_id       VARCHAR(50)  NOT NULL DEFAULT '',
    analyst_name     VARCHAR(150) NOT NULL DEFAULT '',
    pathologist_id   VARCHAR(50)  NOT NULL DEFAULT '',
    pathologist_name VARCHAR(150) NOT NULL DEFAULT '',
    is_verified      BOOLEAN      NOT NULL DEFAULT FALSE,
    overall_notes    TEXT         NOT NULL DEFAULT '',
    verified_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_lab_results_order ON lab_results (order_id);
CREATE INDEX idx_lab_results_verified ON lab_results (is_verified);
