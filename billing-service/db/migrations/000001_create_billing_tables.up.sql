CREATE TABLE tariffs (
    id          CHAR(26)      PRIMARY KEY,
    code        VARCHAR(50)   NOT NULL UNIQUE,
    name        VARCHAR(150)  NOT NULL,
    category    SMALLINT      NOT NULL DEFAULT 1,
    class_id    VARCHAR(50)   NOT NULL DEFAULT 'ALL',
    amount      NUMERIC(15,2) NOT NULL DEFAULT 0,
    description TEXT          NOT NULL DEFAULT '',
    is_active   BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_tariffs_code ON tariffs (code);
CREATE INDEX idx_tariffs_category ON tariffs (category);
CREATE INDEX idx_tariffs_is_active ON tariffs (is_active);

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

CREATE TABLE invoice_items (
    id              CHAR(26)      PRIMARY KEY,
    invoice_id      CHAR(26)      NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    tariff_id       VARCHAR(50)   NOT NULL DEFAULT '',
    item_name       VARCHAR(150)  NOT NULL,
    category        SMALLINT      NOT NULL DEFAULT 1,
    quantity        INT           NOT NULL DEFAULT 1,
    unit_price      NUMERIC(15,2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    subtotal        NUMERIC(15,2) NOT NULL DEFAULT 0,
    notes           TEXT          NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX idx_invoice_items_invoice ON invoice_items (invoice_id);

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

CREATE TABLE payment_items (
    id               CHAR(26)      PRIMARY KEY,
    payment_id       CHAR(26)      NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    method           SMALLINT      NOT NULL DEFAULT 1,
    amount           NUMERIC(15,2) NOT NULL DEFAULT 0,
    reference_number VARCHAR(100)  NOT NULL DEFAULT '',
    notes            TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_payment_items_payment ON payment_items (payment_id);

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

CREATE TABLE deposit_balances (
    patient_id    VARCHAR(50)   PRIMARY KEY,
    total_balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);
