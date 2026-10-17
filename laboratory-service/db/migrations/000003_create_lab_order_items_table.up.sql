CREATE TABLE lab_order_items (
    id        CHAR(26)      PRIMARY KEY,
    order_id  CHAR(26)      NOT NULL REFERENCES lab_orders(id) ON DELETE CASCADE,
    test_id   CHAR(26)      NOT NULL REFERENCES lab_tests(id) ON DELETE RESTRICT,
    test_code VARCHAR(50)   NOT NULL,
    test_name VARCHAR(150)  NOT NULL,
    price     NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    notes     TEXT          NOT NULL DEFAULT ''
);

CREATE INDEX idx_lab_order_items_order ON lab_order_items (order_id);
CREATE INDEX idx_lab_order_items_test ON lab_order_items (test_id);
