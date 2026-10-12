CREATE TABLE beds (
    id         CHAR(26)    PRIMARY KEY,
    code       VARCHAR(50) NOT NULL,
    ward_id    CHAR(26)    NOT NULL REFERENCES wards (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    status     SMALLINT    NOT NULL DEFAULT 1,
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_beds_ward_code ON beds (ward_id, code);
CREATE INDEX idx_beds_ward_id ON beds (ward_id);
CREATE INDEX idx_beds_status ON beds (status);
CREATE INDEX idx_beds_is_active ON beds (is_active);
