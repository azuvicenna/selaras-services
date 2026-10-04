CREATE TABLE patients (
    id                   CHAR(26) PRIMARY KEY,
    medical_record_no    VARCHAR(20)  NOT NULL UNIQUE,
    satusehat_id         VARCHAR(64)  NOT NULL DEFAULT '',
    nik                  VARCHAR(16)  NOT NULL DEFAULT '',
    name                 VARCHAR(255) NOT NULL,
    mother_name          VARCHAR(255) NOT NULL DEFAULT '',

    birth_place          VARCHAR(100) NOT NULL DEFAULT '',
    birth_date           DATE,
    gender               SMALLINT     NOT NULL DEFAULT 0,
    blood_type           SMALLINT     NOT NULL DEFAULT 0,
    marital_status       SMALLINT     NOT NULL DEFAULT 0,
    religion             SMALLINT     NOT NULL DEFAULT 0,

    phone                VARCHAR(20)  NOT NULL DEFAULT '',
    email                VARCHAR(255) NOT NULL DEFAULT '',
    address              TEXT         NOT NULL DEFAULT '',
    village_code         VARCHAR(20)  NOT NULL DEFAULT '',
    district_code        VARCHAR(20)  NOT NULL DEFAULT '',
    city_code            VARCHAR(20)  NOT NULL DEFAULT '',
    province_code        VARCHAR(20)  NOT NULL DEFAULT '',
    postal_code          VARCHAR(10)  NOT NULL DEFAULT '',
    rt                   VARCHAR(5)   NOT NULL DEFAULT '',
    rw                   VARCHAR(5)   NOT NULL DEFAULT '',

    occupation           VARCHAR(100) NOT NULL DEFAULT '',
    education            VARCHAR(100) NOT NULL DEFAULT '',
    preferred_language   VARCHAR(50)  NOT NULL DEFAULT '',

    disability_type      SMALLINT     NOT NULL DEFAULT 0,
    special_needs_note   TEXT         NOT NULL DEFAULT '',

    photo_url            TEXT         NOT NULL DEFAULT '',
    fingerprint_template BYTEA,

    status               SMALLINT     NOT NULL DEFAULT 1,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_patients_nik          ON patients (nik)          WHERE nik <> '';
CREATE UNIQUE INDEX idx_patients_satusehat_id ON patients (satusehat_id) WHERE satusehat_id <> '';
CREATE INDEX        idx_patients_name         ON patients (name);

CREATE SEQUENCE patient_norm_seq;