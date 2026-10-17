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
