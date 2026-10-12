CREATE TABLE icd9cm (
    id          CHAR(26) PRIMARY KEY,
    code        VARCHAR(20) NOT NULL UNIQUE,
    description TEXT        NOT NULL,
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_icd9cm_code ON icd9cm (code);
CREATE INDEX idx_icd9cm_description ON icd9cm (description);
CREATE INDEX idx_icd9cm_is_active ON icd9cm (is_active);
