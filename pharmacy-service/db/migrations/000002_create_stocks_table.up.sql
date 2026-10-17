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
