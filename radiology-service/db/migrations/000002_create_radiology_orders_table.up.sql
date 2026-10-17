CREATE TABLE radiology_orders (
    id                  CHAR(26)      PRIMARY KEY,
    order_number        VARCHAR(50)   NOT NULL UNIQUE,
    encounter_id        VARCHAR(50)   NOT NULL,
    patient_id          VARCHAR(50)   NOT NULL,
    patient_name        VARCHAR(150)  NOT NULL,
    patient_gender      VARCHAR(20)   NOT NULL DEFAULT '',
    patient_dob         VARCHAR(10)   NOT NULL DEFAULT '',
    doctor_id           VARCHAR(50)   NOT NULL DEFAULT '',
    doctor_name         VARCHAR(150)  NOT NULL DEFAULT '',
    clinical_indication TEXT          NOT NULL DEFAULT '',
    priority            SMALLINT      NOT NULL DEFAULT 1,
    status              SMALLINT      NOT NULL DEFAULT 1,
    total_price         NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_radiology_orders_number ON radiology_orders (order_number);
CREATE INDEX idx_radiology_orders_encounter ON radiology_orders (encounter_id);
CREATE INDEX idx_radiology_orders_patient ON radiology_orders (patient_id);
CREATE INDEX idx_radiology_orders_status ON radiology_orders (status);
CREATE INDEX idx_radiology_orders_created_at ON radiology_orders (created_at);
