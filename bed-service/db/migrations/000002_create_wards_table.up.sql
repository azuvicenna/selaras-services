CREATE TABLE wards (
    id            CHAR(26)     PRIMARY KEY,
    code          VARCHAR(50)  NOT NULL UNIQUE,
    name          VARCHAR(100) NOT NULL,
    department_id CHAR(26)     NOT NULL DEFAULT '',
    ward_class_id CHAR(26)     NOT NULL REFERENCES ward_classes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    floor         INT          NOT NULL DEFAULT 1,
    capacity      INT          NOT NULL DEFAULT 1,
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_wards_code ON wards (code);
CREATE INDEX idx_wards_ward_class_id ON wards (ward_class_id);
CREATE INDEX idx_wards_department_id ON wards (department_id) WHERE department_id <> '';
CREATE INDEX idx_wards_floor ON wards (floor);
CREATE INDEX idx_wards_is_active ON wards (is_active);
