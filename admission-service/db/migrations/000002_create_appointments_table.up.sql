CREATE TABLE appointments (
    id               CHAR(26)     PRIMARY KEY,
    booking_code     VARCHAR(50)  NOT NULL UNIQUE,
    patient_id       VARCHAR(50)  NOT NULL,
    practitioner_id  VARCHAR(50)  NOT NULL,
    department_id    VARCHAR(50)  NOT NULL,
    schedule_id      VARCHAR(50)  NOT NULL DEFAULT '',
    appointment_date VARCHAR(10)  NOT NULL,
    queue_number     INT          NOT NULL DEFAULT 1,
    status           SMALLINT     NOT NULL DEFAULT 1,
    notes            TEXT         NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_appointments_booking_code ON appointments (booking_code);
CREATE INDEX idx_appointments_patient ON appointments (patient_id);
CREATE INDEX idx_appointments_practitioner ON appointments (practitioner_id);
CREATE INDEX idx_appointments_department ON appointments (department_id);
CREATE INDEX idx_appointments_date ON appointments (appointment_date);
CREATE INDEX idx_appointments_status ON appointments (status);
CREATE INDEX idx_appointments_schedule ON appointments (schedule_id);
CREATE UNIQUE INDEX idx_appointments_practitioner_date_queue ON appointments (practitioner_id, appointment_date, queue_number);
