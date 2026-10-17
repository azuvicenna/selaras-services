CREATE TABLE deposit_balances (
    patient_id    VARCHAR(50)   PRIMARY KEY,
    total_balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);
