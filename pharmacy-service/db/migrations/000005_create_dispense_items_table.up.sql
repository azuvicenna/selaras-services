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
