CREATE TABLE invoices (
    id              CHAR(26)      PRIMARY KEY,
    invoice_number  VARCHAR(50)   NOT NULL UNIQUE,
    encounter_id    VARCHAR(50)   NOT NULL,
    patient_id      VARCHAR(50)   NOT NULL,
    patient_name    VARCHAR(150)  NOT NULL,
    status          SMALLINT      NOT NULL DEFAULT 1,
    total_amount    NUMERIC(15,2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    paid_amount     NUMERIC(15,2) NOT NULL DEFAULT 0,
    balance_amount  NUMERIC(15,2) NOT NULL DEFAULT 0,
    notes           TEXT          NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_invoices_number ON invoices (invoice_number);
CREATE INDEX idx_invoices_encounter ON invoices (encounter_id);
CREATE INDEX idx_invoices_patient ON invoices (patient_id);
CREATE INDEX idx_invoices_status ON invoices (status);
