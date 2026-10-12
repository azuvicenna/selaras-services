CREATE TABLE icd10 (
    id             CHAR(26) PRIMARY KEY,
    code           VARCHAR(20)  NOT NULL UNIQUE,
    description_id TEXT         NOT NULL,
    description_en TEXT         NOT NULL DEFAULT '',
    category       VARCHAR(100) NOT NULL DEFAULT '',
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_icd10_code ON icd10 (code);
CREATE INDEX idx_icd10_category ON icd10 (category);
CREATE INDEX idx_icd10_description_id ON icd10 (description_id);
CREATE INDEX idx_icd10_is_active ON icd10 (is_active);
