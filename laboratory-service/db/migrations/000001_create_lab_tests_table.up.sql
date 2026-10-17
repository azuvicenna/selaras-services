CREATE TABLE lab_tests (
    id                      CHAR(26)      PRIMARY KEY,
    code                    VARCHAR(50)   NOT NULL UNIQUE,
    name                    VARCHAR(150)  NOT NULL,
    category                SMALLINT      NOT NULL DEFAULT 1,
    loinc_code              VARCHAR(50)   NOT NULL DEFAULT '',
    specimen_type           SMALLINT      NOT NULL DEFAULT 1,
    unit                    VARCHAR(30)   NOT NULL DEFAULT '',
    reference_range_male    VARCHAR(100)  NOT NULL DEFAULT '',
    reference_range_female  VARCHAR(100)  NOT NULL DEFAULT '',
    ref_low                 NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    ref_high                NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    critical_low            NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    critical_high           NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    price                   NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    is_active               BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_lab_tests_code ON lab_tests (code);
CREATE INDEX idx_lab_tests_category ON lab_tests (category);
CREATE INDEX idx_lab_tests_is_active ON lab_tests (is_active);
CREATE INDEX idx_lab_tests_loinc_code ON lab_tests (loinc_code);
