CREATE TABLE tariffs (
    id          CHAR(26)      PRIMARY KEY,
    code        VARCHAR(50)   NOT NULL UNIQUE,
    name        VARCHAR(150)  NOT NULL,
    category    SMALLINT      NOT NULL DEFAULT 1,
    class_id    VARCHAR(50)   NOT NULL DEFAULT 'ALL',
    amount      NUMERIC(15,2) NOT NULL DEFAULT 0,
    description TEXT          NOT NULL DEFAULT '',
    is_active   BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_tariffs_code ON tariffs (code);
CREATE INDEX idx_tariffs_category ON tariffs (category);
CREATE INDEX idx_tariffs_is_active ON tariffs (is_active);
