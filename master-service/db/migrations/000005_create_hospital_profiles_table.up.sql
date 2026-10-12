CREATE TABLE hospital_profiles (
    id                         CHAR(26)     PRIMARY KEY,
    name                       VARCHAR(255) NOT NULL,
    bpjs_facility_code         VARCHAR(50)  NOT NULL DEFAULT '',
    kemenkes_facility_code     VARCHAR(50)  NOT NULL DEFAULT '',
    tax_identification_number VARCHAR(50)  NOT NULL DEFAULT '',
    address                    TEXT         NOT NULL DEFAULT '',
    contact_email              VARCHAR(100) NOT NULL DEFAULT '',
    contact_phone              VARCHAR(50)  NOT NULL DEFAULT '',
    created_at                 TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at                 TIMESTAMPTZ  NOT NULL DEFAULT now()
);
