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
