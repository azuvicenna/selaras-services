CREATE TABLE lab_specimens (
    id              CHAR(26)     PRIMARY KEY,
    order_id        CHAR(26)     NOT NULL REFERENCES lab_orders(id) ON DELETE CASCADE,
    specimen_number VARCHAR(50)  NOT NULL UNIQUE,
    specimen_type   SMALLINT     NOT NULL DEFAULT 1,
    status          SMALLINT     NOT NULL DEFAULT 1,
    collector_id    VARCHAR(50)  NOT NULL DEFAULT '',
    collector_name  VARCHAR(150) NOT NULL DEFAULT '',
    notes           TEXT         NOT NULL DEFAULT '',
    collected_at    TIMESTAMPTZ,
    received_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_lab_specimens_number ON lab_specimens (specimen_number);
CREATE INDEX idx_lab_specimens_order ON lab_specimens (order_id);
CREATE INDEX idx_lab_specimens_status ON lab_specimens (status);
