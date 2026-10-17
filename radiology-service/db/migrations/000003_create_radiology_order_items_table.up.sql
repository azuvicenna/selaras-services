CREATE TABLE radiology_order_items (
    id         CHAR(26)      PRIMARY KEY,
    order_id   CHAR(26)      NOT NULL REFERENCES radiology_orders(id) ON DELETE CASCADE,
    exam_id    CHAR(26)      NOT NULL REFERENCES radiology_exams(id) ON DELETE RESTRICT,
    exam_code  VARCHAR(50)   NOT NULL,
    exam_name  VARCHAR(150)  NOT NULL,
    modality   SMALLINT      NOT NULL DEFAULT 1,
    price      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    notes      TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_radiology_order_items_order ON radiology_order_items (order_id);
CREATE INDEX idx_radiology_order_items_exam ON radiology_order_items (exam_id);
