CREATE TABLE lab_result_items (
    id              CHAR(26)      PRIMARY KEY,
    result_id       CHAR(26)      NOT NULL REFERENCES lab_results(id) ON DELETE CASCADE,
    order_id        CHAR(26)      NOT NULL REFERENCES lab_orders(id) ON DELETE CASCADE,
    test_id         CHAR(26)      NOT NULL REFERENCES lab_tests(id) ON DELETE RESTRICT,
    test_name       VARCHAR(150)  NOT NULL,
    loinc_code      VARCHAR(50)   NOT NULL DEFAULT '',
    result_value    VARCHAR(100)  NOT NULL DEFAULT '',
    numeric_value   NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    reference_range VARCHAR(100)  NOT NULL DEFAULT '',
    unit            VARCHAR(30)   NOT NULL DEFAULT '',
    flag            SMALLINT      NOT NULL DEFAULT 1,
    notes           TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_lab_result_items_result ON lab_result_items (result_id);
CREATE INDEX idx_lab_result_items_order ON lab_result_items (order_id);
