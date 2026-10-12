CREATE TABLE departments (
    id              CHAR(26)    PRIMARY KEY,
    code            VARCHAR(50)  NOT NULL UNIQUE,
    name            VARCHAR(100) NOT NULL,
    department_type SMALLINT     NOT NULL DEFAULT 0,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_departments_code ON departments (code);
CREATE INDEX idx_departments_department_type ON departments (department_type);
CREATE INDEX idx_departments_is_active ON departments (is_active);
