CREATE TABLE medications (
    id                CHAR(26)      PRIMARY KEY,
    kfa_code          VARCHAR(64)   NOT NULL UNIQUE,
    name              VARCHAR(255)  NOT NULL,
    form              INT           NOT NULL DEFAULT 1,
    strength          VARCHAR(64)   NOT NULL DEFAULT '',
    unit              VARCHAR(32)   NOT NULL DEFAULT '',
    package_unit      VARCHAR(32)   NOT NULL DEFAULT '',
    conversion_factor INT           NOT NULL DEFAULT 1,
    selling_price     NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    purchase_price    NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    min_stock         INT           NOT NULL DEFAULT 0,
    is_active         BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_medications_kfa_code ON medications (kfa_code);
CREATE INDEX idx_medications_name ON medications (name);
CREATE INDEX idx_medications_form ON medications (form);
CREATE INDEX idx_medications_is_active ON medications (is_active);

CREATE TABLE stocks (
    id               CHAR(26)    PRIMARY KEY,
    medication_id    CHAR(26)    NOT NULL REFERENCES medications(id) ON DELETE CASCADE,
    batch_number     VARCHAR(64) NOT NULL,
    expiry_date      VARCHAR(10) NOT NULL,
    current_quantity INT         NOT NULL DEFAULT 0,
    location         VARCHAR(64) NOT NULL DEFAULT 'Gudang Utama',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_stocks_med_batch_loc UNIQUE (medication_id, batch_number, location)
);

CREATE INDEX idx_stocks_medication_id ON stocks (medication_id);
CREATE INDEX idx_stocks_location ON stocks (location);
CREATE INDEX idx_stocks_expiry_date ON stocks (expiry_date);

CREATE TABLE stock_movements (
    id             CHAR(26)    PRIMARY KEY,
    medication_id  CHAR(26)    NOT NULL REFERENCES medications(id) ON DELETE CASCADE,
    batch_number   VARCHAR(64) NOT NULL,
    movement_type  INT         NOT NULL,
    quantity       INT         NOT NULL,
    reference_type VARCHAR(64) NOT NULL DEFAULT '',
    reference_id   VARCHAR(64) NOT NULL DEFAULT '',
    notes          TEXT        NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_stock_movements_medication_id ON stock_movements (medication_id);
CREATE INDEX idx_stock_movements_type ON stock_movements (movement_type);
CREATE INDEX idx_stock_movements_created_at ON stock_movements (created_at);

CREATE TABLE dispenses (
    id              CHAR(26)      PRIMARY KEY,
    dispense_number VARCHAR(64)   NOT NULL UNIQUE,
    prescription_id VARCHAR(64)   NOT NULL,
    encounter_id    VARCHAR(64)   NOT NULL,
    patient_id      VARCHAR(64)   NOT NULL,
    patient_name    VARCHAR(255)  NOT NULL,
    pharmacist_id   VARCHAR(64)   NOT NULL DEFAULT '',
    status          INT           NOT NULL DEFAULT 1,
    total_amount    NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    notes           TEXT          NOT NULL DEFAULT '',
    dispensed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_dispenses_number ON dispenses (dispense_number);
CREATE INDEX idx_dispenses_prescription_id ON dispenses (prescription_id);
CREATE INDEX idx_dispenses_encounter_id ON dispenses (encounter_id);
CREATE INDEX idx_dispenses_patient_id ON dispenses (patient_id);
CREATE INDEX idx_dispenses_status ON dispenses (status);

CREATE TABLE dispense_items (
    id                 CHAR(26)      PRIMARY KEY,
    dispense_id        CHAR(26)      NOT NULL REFERENCES dispenses(id) ON DELETE CASCADE,
    medication_id      CHAR(26)      NOT NULL REFERENCES medications(id) ON DELETE RESTRICT,
    medication_name    VARCHAR(255)  NOT NULL,
    batch_number       VARCHAR(64)   NOT NULL DEFAULT '',
    quantity           INT           NOT NULL DEFAULT 1,
    unit_price         NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    subtotal           NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    dosage_instruction TEXT          NOT NULL DEFAULT '',
    notes              TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_dispense_items_dispense_id ON dispense_items (dispense_id);
CREATE INDEX idx_dispense_items_medication_id ON dispense_items (medication_id);
