CREATE TABLE ward_classes (
    id                    CHAR(26)     PRIMARY KEY,
    code                  VARCHAR(20)  NOT NULL UNIQUE,
    name                  VARCHAR(100) NOT NULL,
    bpjs_equivalent_class VARCHAR(50)  NOT NULL DEFAULT '',
    is_active             BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_ward_classes_code ON ward_classes (code);
CREATE INDEX idx_ward_classes_is_active ON ward_classes (is_active);
