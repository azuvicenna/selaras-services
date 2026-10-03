CREATE TABLE insurances (
    id             CHAR(26) PRIMARY KEY,
    patient_id     CHAR(26)     NOT NULL REFERENCES patients (id),
    provider       SMALLINT     NOT NULL DEFAULT 0,
    policy_number  VARCHAR(50)  NOT NULL,
    insurance_name VARCHAR(100) NOT NULL DEFAULT '',
    class_type     VARCHAR(50)  NOT NULL DEFAULT '',
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    UNIQUE (patient_id, policy_number)
);