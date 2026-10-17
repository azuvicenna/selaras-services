CREATE TABLE payment_items (
    id               CHAR(26)      PRIMARY KEY,
    payment_id       CHAR(26)      NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    method           SMALLINT      NOT NULL DEFAULT 1,
    amount           NUMERIC(15,2) NOT NULL DEFAULT 0,
    reference_number VARCHAR(100)  NOT NULL DEFAULT '',
    notes            TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_payment_items_payment ON payment_items (payment_id);
