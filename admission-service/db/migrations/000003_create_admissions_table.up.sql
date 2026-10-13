CREATE TABLE admissions (
    id               CHAR(26)     PRIMARY KEY,
    admission_number VARCHAR(50)  NOT NULL UNIQUE,
    patient_id       VARCHAR(50)  NOT NULL,
    appointment_id   VARCHAR(50),
    admission_type   SMALLINT     NOT NULL DEFAULT 1,
    status           SMALLINT     NOT NULL DEFAULT 1,
    department_id    VARCHAR(50)  NOT NULL,
    practitioner_id  VARCHAR(50)  NOT NULL,
    admission_date   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    discharge_date   TIMESTAMPTZ,
    ward_id          VARCHAR(50),
    bed_id           VARCHAR(50),
    notes            TEXT         NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_admissions_number ON admissions (admission_number);
CREATE INDEX idx_admissions_patient ON admissions (patient_id);
CREATE INDEX idx_admissions_type ON admissions (admission_type);
CREATE INDEX idx_admissions_status ON admissions (status);
CREATE INDEX idx_admissions_department ON admissions (department_id);
CREATE INDEX idx_admissions_practitioner ON admissions (practitioner_id);
CREATE INDEX idx_admissions_date ON admissions (admission_date);
