CREATE TABLE schedules (
    id              CHAR(26)     PRIMARY KEY,
    practitioner_id VARCHAR(50)  NOT NULL,
    department_id   VARCHAR(50)  NOT NULL,
    day_of_week     SMALLINT     NOT NULL,
    start_time      VARCHAR(10)  NOT NULL,
    end_time        VARCHAR(10)  NOT NULL,
    quota           INT          NOT NULL DEFAULT 10,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_schedules_practitioner ON schedules (practitioner_id);
CREATE INDEX idx_schedules_department ON schedules (department_id);
CREATE INDEX idx_schedules_day_of_week ON schedules (day_of_week);
CREATE INDEX idx_schedules_is_active ON schedules (is_active);
