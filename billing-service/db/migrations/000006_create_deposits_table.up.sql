CREATE TABLE deposits (
    id               CHAR(26)      PRIMARY KEY,
    patient_id       VARCHAR(50)   NOT NULL,
    encounter_id     VARCHAR(50)   NOT NULL DEFAULT '',
    deposit_type     SMALLINT      NOT NULL DEFAULT 1,
    amount           NUMERIC(15,2) NOT NULL DEFAULT 0,
    payment_method   SMALLINT      NOT NULL DEFAULT 1,
    reference_number VARCHAR(100)  NOT NULL DEFAULT '',
    notes            TEXT          NOT NULL DEFAULT '',
    cashier_id       VARCHAR(50)   NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX idx_deposits_patient ON deposits (patient_id);
