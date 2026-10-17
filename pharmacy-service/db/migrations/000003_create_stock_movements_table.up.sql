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
