CREATE TABLE medications (
    id                CHAR(26)      PRIMARY KEY,
    kfa_code          VARCHAR(64)   NOT NULL UNIQUE,
    name              VARCHAR(255)  NOT NULL,
    form              INT           NOT NULL DEFAULT 1,
    strength          VARCHAR(64)   NOT NULL DEFAULT '',
    unit              VARCHAR(32)   NOT NULL DEFAULT '',
    package_unit      VARCHAR(32)   NOT NULL DEFAULT '',
    conversion_factor INT           NOT NULL DEFAULT 1,
    selling_price     NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    purchase_price    NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    min_stock         INT           NOT NULL DEFAULT 0,
    is_active         BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_medications_kfa_code ON medications (kfa_code);
CREATE INDEX idx_medications_name ON medications (name);
CREATE INDEX idx_medications_form ON medications (form);
CREATE INDEX idx_medications_is_active ON medications (is_active);
