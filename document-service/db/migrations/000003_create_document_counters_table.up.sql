CREATE TABLE document_counters (
    prefix           VARCHAR(20) NOT NULL,
    year             INT         NOT NULL,
    current_sequence BIGINT      NOT NULL DEFAULT 0,
    PRIMARY KEY (prefix, year)
);
