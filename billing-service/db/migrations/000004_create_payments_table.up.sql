CREATE TABLE payments (
    id             CHAR(26)      PRIMARY KEY,
    receipt_number VARCHAR(50)   NOT NULL UNIQUE,
    invoice_id     CHAR(26)      NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    cashier_id     VARCHAR(50)   NOT NULL DEFAULT '',
    cashier_name   VARCHAR(150)  NOT NULL DEFAULT '',
    total_paid     NUMERIC(15,2) NOT NULL DEFAULT 0,
    paid_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_payments_receipt ON payments (receipt_number);
CREATE INDEX idx_payments_invoice ON payments (invoice_id);
