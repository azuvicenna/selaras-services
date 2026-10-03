CREATE SEQUENCE patient_medical_record_seq;

CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_record_no TEXT NOT NULL UNIQUE
        DEFAULT ('RM-' || lpad(nextval('patient_medical_record_seq')::text, 8, '0')),
    nik TEXT NOT NULL UNIQUE CHECK (nik ~ '^\d{16}$'),
    name TEXT NOT NULL,
    birth_date DATE NOT NULL,
    gender TEXT NOT NULL CHECK (gender IN ('male', 'female')),
    phone TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);