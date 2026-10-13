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

CREATE TABLE lab_orders (
    id             CHAR(26)      PRIMARY KEY,
    order_number   VARCHAR(50)   NOT NULL UNIQUE,
    encounter_id   VARCHAR(50)   NOT NULL,
    patient_id     VARCHAR(50)   NOT NULL,
    patient_name   VARCHAR(150)  NOT NULL,
    patient_gender VARCHAR(20)   NOT NULL DEFAULT '',
    patient_dob    VARCHAR(10)   NOT NULL DEFAULT '',
    doctor_id      VARCHAR(50)   NOT NULL DEFAULT '',
    doctor_name    VARCHAR(150)  NOT NULL DEFAULT '',
    clinical_notes TEXT          NOT NULL DEFAULT '',
    priority       SMALLINT      NOT NULL DEFAULT 1,
    status         SMALLINT      NOT NULL DEFAULT 1,
    total_price    NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_lab_orders_number ON lab_orders (order_number);
CREATE INDEX idx_lab_orders_encounter ON lab_orders (encounter_id);
CREATE INDEX idx_lab_orders_patient ON lab_orders (patient_id);
CREATE INDEX idx_lab_orders_status ON lab_orders (status);
CREATE INDEX idx_lab_orders_created_at ON lab_orders (created_at);

CREATE TABLE lab_order_items (
    id        CHAR(26)      PRIMARY KEY,
    order_id  CHAR(26)      NOT NULL REFERENCES lab_orders(id) ON DELETE CASCADE,
    test_id   CHAR(26)      NOT NULL REFERENCES lab_tests(id) ON DELETE RESTRICT,
    test_code VARCHAR(50)   NOT NULL,
    test_name VARCHAR(150)  NOT NULL,
    price     NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    notes     TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_lab_order_items_order ON lab_order_items (order_id);
CREATE INDEX idx_lab_order_items_test ON lab_order_items (test_id);

CREATE TABLE lab_specimens (
    id              CHAR(26)     PRIMARY KEY,
    order_id        CHAR(26)     NOT NULL REFERENCES lab_orders(id) ON DELETE CASCADE,
    specimen_number VARCHAR(50)  NOT NULL UNIQUE,
    specimen_type   SMALLINT     NOT NULL DEFAULT 1,
    status          SMALLINT     NOT NULL DEFAULT 1,
    collector_id    VARCHAR(50)  NOT NULL DEFAULT '',
    collector_name  VARCHAR(150) NOT NULL DEFAULT '',
    notes           TEXT         NOT NULL DEFAULT '',
    collected_at    TIMESTAMPTZ,
    received_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_lab_specimens_number ON lab_specimens (specimen_number);
CREATE INDEX idx_lab_specimens_order ON lab_specimens (order_id);
CREATE INDEX idx_lab_specimens_status ON lab_specimens (status);

CREATE TABLE lab_results (
    id              CHAR(26)     PRIMARY KEY,
    order_id        CHAR(26)     NOT NULL UNIQUE REFERENCES lab_orders(id) ON DELETE CASCADE,
    order_number    VARCHAR(50)  NOT NULL,
    patient_id      VARCHAR(50)  NOT NULL,
    analyst_id      VARCHAR(50)  NOT NULL DEFAULT '',
    analyst_name    VARCHAR(150) NOT NULL DEFAULT '',
    pathologist_id  VARCHAR(50)  NOT NULL DEFAULT '',
    pathologist_name VARCHAR(150) NOT NULL DEFAULT '',
    is_verified     BOOLEAN      NOT NULL DEFAULT FALSE,
    overall_notes   TEXT         NOT NULL DEFAULT '',
    verified_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_lab_results_order ON lab_results (order_id);
CREATE INDEX idx_lab_results_verified ON lab_results (is_verified);

CREATE TABLE lab_result_items (
    id              CHAR(26)      PRIMARY KEY,
    result_id       CHAR(26)      NOT NULL REFERENCES lab_results(id) ON DELETE CASCADE,
    order_id        CHAR(26)      NOT NULL REFERENCES lab_orders(id) ON DELETE CASCADE,
    test_id         CHAR(26)      NOT NULL REFERENCES lab_tests(id) ON DELETE RESTRICT,
    test_name       VARCHAR(150)  NOT NULL,
    loinc_code      VARCHAR(50)   NOT NULL DEFAULT '',
    result_value    VARCHAR(100)  NOT NULL DEFAULT '',
    numeric_value   NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    reference_range VARCHAR(100)  NOT NULL DEFAULT '',
    unit            VARCHAR(30)   NOT NULL DEFAULT '',
    flag            SMALLINT      NOT NULL DEFAULT 1,
    notes           TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_lab_result_items_result ON lab_result_items (result_id);
CREATE INDEX idx_lab_result_items_order ON lab_result_items (order_id);
