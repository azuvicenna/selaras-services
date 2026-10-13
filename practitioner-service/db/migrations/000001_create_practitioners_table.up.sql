CREATE TABLE practitioners (
    id                CHAR(26)     PRIMARY KEY,
    nik               VARCHAR(20)  NOT NULL UNIQUE,
    name              VARCHAR(150) NOT NULL,
    gender            SMALLINT     NOT NULL DEFAULT 0,
    birth_date        VARCHAR(10)  NOT NULL DEFAULT '',
    phone             VARCHAR(20)  NOT NULL DEFAULT '',
    email             VARCHAR(100) NOT NULL DEFAULT '',
    practitioner_type SMALLINT     NOT NULL DEFAULT 1,
    specialization    VARCHAR(100) NOT NULL DEFAULT '',
    sip_number        VARCHAR(100) NOT NULL DEFAULT '',
    sip_expiry_date   VARCHAR(10)  NOT NULL DEFAULT '',
    str_number        VARCHAR(100) NOT NULL DEFAULT '',
    str_expiry_date   VARCHAR(10)  NOT NULL DEFAULT '',
    department_id     VARCHAR(50)  NOT NULL DEFAULT '',
    ihs_number        VARCHAR(50)  NOT NULL DEFAULT '',
    is_active         BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_practitioners_nik ON practitioners (nik);
CREATE INDEX idx_practitioners_sip ON practitioners (sip_number);
CREATE INDEX idx_practitioners_type ON practitioners (practitioner_type);
CREATE INDEX idx_practitioners_department ON practitioners (department_id);
CREATE INDEX idx_practitioners_is_active ON practitioners (is_active);
